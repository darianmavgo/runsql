# runsql
Golang CLI run a list of sql files

# tried and failed
I wanted to see if I could faster better results by not using readutil to get all the sql but
- Things like ".read myscript.sql " can't be executed by db.Exec.   No dot commands work with database/sql
- Using os/exec to run commands ended up not being worth the effort to troubleshoot. 
