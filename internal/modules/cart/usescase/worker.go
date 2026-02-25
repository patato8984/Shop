package cart_usescase

import (
	"context"
	"fmt"
	"sync"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
	cart_repo "github.com/patato8984/Shop/internal/modules/cart/repo"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
)

type ProductPriceProvider interface {
	WorkerGetPrice(ctx context.Context, idSkus int) (float64, error)
}
type CacheDel interface {
	Del(ctx context.Context, key string) bool
}
type PriceUpdaterService struct {
	provider ProductPriceProvider
	repoCart cart_repo.CartRepo
	kp       shared_events.EventPublisher
}

func NewWorkerService(provider ProductPriceProvider, repoCart *cart_repo.CartRepo, kp shared_events.EventPublisher) *PriceUpdaterService {
	return &PriceUpdaterService{provider: provider, repoCart: *repoCart, kp: kp}
}
func (r *PriceUpdaterService) UpdateWorkerPrice(ctx context.Context) error {
	price, err := r.repoCart.WorkerGetAllProductId(ctx)
	if err != nil {
		return err
	}
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)
	errorDB := make(chan error, 100)
OuterLoop:
	for id_skus, id_cart := range price {
		select {
		case <-workerCtx.Done():
			break OuterLoop
		default:
		}
		wg.Add(1)
		semaphore <- struct{}{}
		go func(ctx context.Context, sID int, cIDs []int, errorDB chan error) {
			select {
			case <-workerCtx.Done():
				return
			default:
			}
			defer func() {
				<-semaphore
				wg.Done()
			}()
			price, err := r.provider.WorkerGetPrice(ctx, sID)
			if err != nil {
				errorDB <- err
				cancel()
				return
			}
			var wg1 sync.WaitGroup
			semaphore1 := make(chan struct{}, 5)
		InnerLoop:
			for _, cartID := range cIDs {
				select {
				case <-workerCtx.Done():
					break InnerLoop
				default:
				}
				semaphore1 <- struct{}{}
				wg1.Add(1)
				go func(ctx context.Context, cartID int, price float64, sID int, errorDB chan error) {
					defer func() {
						<-semaphore1
						wg1.Done()
					}()
					err := r.repoCart.WorkerUpdatePrice(ctx, cartID, sID, price)
					if err != nil {
						errorDB <- err
						cancel()
						return
					}
					event := cart_model.CartEvent{
						EventType: "cartUpdateWorker",
						PayLoad: cart_model.CartUpdate{
							IdCart: cartID,
						},
					}
					if err := r.kp.Publisher(ctx, "cart_event", event.EventType, event); err != nil {
						fmt.Print("error publisher event")
					}
				}(ctx, cartID, price, sID, errorDB)
			}
			wg1.Wait()
		}(ctx, id_skus, id_cart, errorDB)
	}
	go func() {
		wg.Wait()
		close(errorDB)
	}()

	for err := range errorDB {
		if err != nil {
			return err
		}
	}
	return nil
}
