# PgOra-Comparator
Postgres and Oracle Table Data Comparator
Steps to Use :
1. Create a DB user on both Oracle and Postgres with read only access to required tables.
2. Add 2 files with Oracle and PG database details and a file with AWS cognito detials for login acesss.(The handler.go can be modified to skip AWS cognito login process).
3. Add env variables with file path of step2, PG env var nae -pgsecret,Ora env var name - orasecret, App env var name - appconfig
4. Add the Query to compare in the CompQuery.json file.
5. Run the Binary and aceess the application URL to view and compare tables.
