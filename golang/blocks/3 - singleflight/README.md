# Singleflight

Base : Thundering herd problem.

Story : Let's say you have a long-running query or expensive 3rd party API call. And your environment is concurrent.

Goal : suppress the duplicate function call.

Solution : wrap your expensive/long-running function with singleflight! This might decrease latency also!