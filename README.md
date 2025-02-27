# Meeting Scheduler API

A minimalist, horizontally scalable service to solve the time zone issue for distributed teams.

## Technical Architecture

```
events/
  ├── creation → minimal validation, direct MongoDB persistence
  ├── availability → user time slots with IANA timezone encoding
  └── recommendations → algorithm prioritizes max participant overlap
```

## Engineering Decisions & Trade-offs

### Time Zone Handling
- All conversions happen at API boundaries; core logic operates on UTC
- Avoided timezone libraries that add dependency complexity

### Data Model Simplicity
- Events as the core entity; no separate user model to reduce JOIN complexity
- Document-based storage chosen to match the natural hierarchy of our data
- Clean separation between time slot representation and business logic

### Performance Considerations
- Read-heavy workload optimized with MongoDB (document-based retrieval)
- Horizontally scalable API servers (stateless design)
- Query patterns kept simple for predictable performance

### Testing Approach
- Algorithm unit-tested separately from database interactions
- Integration tests with mocked DB using custom Router wrapper
- Pragmatic test coverage focused on recommendation algorithm correctness

### Scalability Limits & Improvements
- Current design handles thousands of concurrent users 
- Identified bottleneck: recommendation calculation with many attendees
- Future: Adding read replicas or implementing algorithmic optimizations

## API Contract

**Core Endpoints:**
```
DELTE/POST/PUT      /events/{id}                            → Create/update/delete events
GET                 /events/{id}                            → Retrieve event details
DELTE/POST/PUT      /events/{id}/availability/{user_id}     → Manage availability
GET                 /events/{id}/recommendations            → Calculate optimal slots
```

## Deployment Architecture

Simple two-container Kubernetes deployment:
1. Stateless API servers (horizontally scalable)
2. MongoDB with persistent storage

``` ascii
       ┌────────────────────────┐             
       │      load balancer     │             
     ┌─┴────────────┬───────────┴───┐         
     │              │               │         
     │              │               │         
     │              │               │         
┌────▼─────┐   ┌────▼─────┐  ┌──────▼───┐     
│api server│   │api server│  │api server│     
│  pod     │   │  pod     │  │  pod     │     
└────┬─────┘   └─────┬────┘  └──────┬───┘     
     │               │              │         
     │               │              │         
     │          ┌────▼─────┐        │         
     └──────────► mongodb  ◄────────┘         
                │   pod    │                  
                └─────┬────┘                  
                      │                       
                      │                       
                      ▼                       
               ┌───────────────┐              
               │  persistent   │              
               │     volume    │              
               └───────────────┘              
```

## What I Would Improve Given More Time

1. Transaction support for concurrent availability updates
2. Proper abstraction layer between handlers and data access
3. Caching layer for frequently-accessed events
4. Distributed tracing for API performance monitoring
5. More robust error handling beyond HTTP status codes
6. Test k8s deployment more thoroughly

## Local Development

```bash
docker-compose up
```

## Cloud Deployment

```bash
kubectl apply -f k8s-simple-mongo.yaml
kubectl apply -f k8s-simple-app.yaml
```

## Local Testing

```bash
go test -v .
```
