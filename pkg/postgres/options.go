package postgres

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func MaxConnAttempts(attempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(d time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = d
	}
}
