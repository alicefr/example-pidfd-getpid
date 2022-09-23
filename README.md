# Example 
```bash
docker run --name unprivileged -ti getfd /usr/bin/proxy
docker inspect unprivileged|grep -i pid
"Pid": 42999,
docker run -ti -d --name privileged \
  --pid host \
  getfd \
  socat UNIX-LISTEN:/tmp/priv.sock -
docker exec -ti privileged connector -pid=42999 -fd=3
```
