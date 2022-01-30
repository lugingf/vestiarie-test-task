# Test task

## To launch

Project use docker-compose. Please use following command to run
`make run`
Or `make rund` to run as a daemon (then logs should be requested by `docker logs`)

Use `make stop` to shutdown

There are two containers. `app` with Go application and `database` with MySQL DB

Request example

```
curl --request POST \
  --url http://localhost:8080/payouts \
  --header 'PaymentUpdateID: 4' \
  --header 'content-type: application/json' \
  --data '[
	{
		"name": "name_1",
		"currency": "EUR",
		"price": 50,
		"seller_id": 1
	},
	{
		"name": "name_2",
		"currency": "EUR",
		"price": 52,
		"seller_id": 1
	},
	{
		"name": "name_3",
		"currency": "USD",
		"price": 1000,
		"seller_id": 2
	},
	{
		"name": "name_4",
		"currency": "USD",
		"price": 10,
		"seller_id": 2
	},
		{
		"name": "name_4",
		"currency": "USD",
		"price": 4000,
		"seller_id": 2
	},
			{
		"name": "name_4",
		"currency": "USD",
		"price": 6000,
		"seller_id": 2
	}
]'
```

## Idempotency

Use header `PaymentUpdateID` to divide updates. If you try to update data with elready existed PaymentUpdateID you'll
got an error

## DataBase

```
Driver:   "mysql",
User:     "gotest",
Password: "gotest",
Host:     "database",
Port:     "3306",
Name:     "local_gotest",
```

### Migrations
Migrations runs automatically with service start. I don't use this approach usually

```
	  CREATE TABLE IF NOT EXISTS item (
	    id bigint PRIMARY KEY AUTO_INCREMENT,
		update_id varchar(16),
		item_name varchar(256),
		price float,
		currency varchar(4),
		seller_id integer,
		INDEX ix_update_id(update_id)
	  );
	  
	  CREATE TABLE IF NOT EXISTS payout (
	    id bigint PRIMARY KEY AUTO_INCREMENT,
		update_id varchar(16),
		seller_id integer,
		amount float,
		currency varchar(4),
		item_id_list text,
		payout_part integer,
	    UNIQUE INDEX ux_sel_cur_upd(seller_id, currency, update_id)
	  );
```

## Payout amount should not exceed a certain limit
Currently it is 5000, and hardcoded. 

## Payout should be linked with at least one Item
I've used lazy approach and just listed items ID for payout. I know it is bad solution. Better would be store `payout_id` 
in `item` table, or create `payout2item` as one-to-many. It is not a big deal, but I've spent a lot of time 
on creating service, making all project works etc. If it is important for the test task we can discuss pros and cons of it.

## Tests

Sorry, but no tests. I've spent enough of time on test task.