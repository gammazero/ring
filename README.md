# ring
Fixed-size circular deque

## Description

Ring is a fixed-size circular queue, where items can be effiifiently added and removed from either end. The number of items that can be stored in the ring is set at creation, and if the capacity is exceeded the end of the Ring being added to overwrites the other end of the Ring.

## Generics

Ring uses generics to create a Ring that contains items of the type specified. To create a Ring that holds a specific type, provide a type argument to `New`. For example:
```go
    stringRing := ring.New[string](10)
```
