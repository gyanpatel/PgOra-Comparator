# PgOra-Comparator
Postgres and Oracle Table Data Comparator
Steps to Use :
1. Create a DB user on both Oracle and Postgres with read only access to required tables.
2. Add 2 files with Oracle and PG database details and a file with AWS cognito detials for login acesss.(The handler.go can be modified ) to skip this.
3. Add env pariables with file path of step2
4. Add the Query to compare in the CompQuery.json file.
5. Run the Binary and aceess the application URL to view and compare tables.
