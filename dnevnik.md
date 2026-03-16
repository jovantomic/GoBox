# 16. mart 2025 dan 1 dev
GO deluje kao kul jezek
- uspeo sam da pokrenem bash

ubuntu@gobox-vm:~/gobox$ ps
    PID TTY          TIME CMD
  11229 pts/1    00:00:00 bash
  11293 pts/1    00:00:00 go
  11329 pts/1    00:00:00 main
  11333 pts/1    00:00:00 bash
  11351 pts/1    00:00:00 ps

  ovo mi daje sve procese u terminali, ideja je da ja mogu da za samo svoj gobox uzmem i narpavim da gobox vidi samo svoje procese i da oni krenu od pid1 i to se radi sa namespacovima