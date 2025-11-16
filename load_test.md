| METHOD | API endpoint            | QPS    | Total requests | Concurrency | Total time (ms)  | Avg per request (ms)   | Used to test                                                 |
|--------|-------------------------|--------|----------------|-------------|------------------|------------------------|--------------------------------------------------------------|
| GET    | /team/get               | 14744  | 1000           | 20          | 67.8             | 1.3                    | hey command                                                  |
| POST   | /team/create            | 16     | 100            | 5           | 6067             | 60.67                  | bash script with random team_names and user_ids (2–5 users)  |
| GET    | /users/getReview        | 10161  | 1000           | 20          | 98.4             | 1.9                    | hey command                                                  |
| POST   | /users/setIsActive      | 2799   | 1000           | 20          | 357.2            | 7                      | hey command                                                  |
| POST   | /pullRequest/merge      | 6924   | 1000           | 20          | 144.4            | 2.8                    | hey command                                                  |
| POST   | /pullRequest/reassign   | 10596  | 1000           | 20          | 94.4             | 1.8                    | hey command                                                  |
| POST   | /pullRequest/create     | 50     | 100            | 5           | 2208             | 22.08                  | bash script with random PR ids                               |
