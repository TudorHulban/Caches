# Memory Caches  

## Cache Vanilla  

Simple memory cache implemented with a Go map and guarded by a mutex.  
Just one file with no dependencies that can be used wherever needed.  
Tested for race conditions.  

## Cache TTL

Added a janitor to clean TTL outdated records.

## LRU Any Version

Could be used when mixed values are to be stored.  
Small penalty in performance.

## LRU Specialised Version

One type of value to be stored.  
Best performance.
