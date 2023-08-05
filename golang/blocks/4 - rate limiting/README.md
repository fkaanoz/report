# Rate Limiting

Goal : limit the api calls with token bucket algorithm.

With global rate limiting middleware, you can limit any API calls regardless of ip addresses. Bucket with 8 token capacity will be filled 4 token per second.

Client based rate limiting middleware will track ip addresses of users in a map, and rate limiting logic will be implemented according to ip addresses. This map can grow infinitely, so we should clean up it.  