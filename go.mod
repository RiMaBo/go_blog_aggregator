module github.com/RiMaBo/go_blog_aggregator

go 1.26.1

replace internal/config => ./internal/config
require internal/config v0.0.0

replace internal/database => ./internal/database
require internal/database v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.12.3 // indirect
)
