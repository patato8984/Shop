logs:
    docker logs -f main4-app-1 | jq -R 'fromjson? | "\(.t | strflocaltime("%Y-%m-%d %H:%M:%S"))  |  ID: \(.request_id // "N/A")  |   LEVEL: \(.l)  |  \(.method)  \(.status) | \(.m)  | \(.url) error: \(.error // "N/A")"' 
logs response:
     docker logs -f --tail 20 main4-app-1 | jq -R 'fromjson? | select(.m == "response")  |  "\(.t | strflocaltime("%Y-%m-%d %H:%M:%S")) | \(.method)  |  \(.status)  |  \(.url)"'
logs error or warn:
    docker logs -f --tail 20 main4-app-1 | jq -R 'fromjson? | select(.l == "warn" or .l == "error")  | "\(.t | strflocaltime("%Y-%m-%d %H:%M:%S"))  |  \(.l)  |  \(.error )" '