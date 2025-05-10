# joshfiles
> makefiles for those who dont want makefiles  

build with `go build . -o ./josh` and put somewhere in your path  
or you can run `josh install` but how would you do it without the thing  

you run `josh --init` and it makes a joshfile in your directory  

## features  
* async execution  
* pretty error reports  
* yaml config  
* reports on time elapsed  
* execute stuff with or without a shell  
* easy to modify  

## plans
* toml support (in case you want both josh and tom on your system)  
* json support (in case you want both josh and jason on your system)  
they're very easy, i just wanna do it tomorrow  

## usage

- **josh**  
  runs the first task defined in the joshfile.

- **josh *TASK***  
  runs a specific named task.

- **josh --init**  
  creates an example configuration file in the current directory.

- **josh -h, --h**  
  displays help information.

- **josh --man**  
  displays this man page.

## example

```yaml
tasks:
  main:
    cmds: # runs stuff in sh
      - "echo hello && echo world"
  direct:
    direct: true # bypasses sh
    cmds:
      - "echo hello && echo world"
  async:
    async: true # runs them concurrently. remember that they'll finish out of order
    report: true # reports time elapsed
    cmds:
      - "echo notice how they all finish at the same time?"
      - "sleep 2"
      - "sleep 2"
      - "sleep 2"
  silent:
    silent: true # runs them silently 
    cmds:
      - "echo you wont see this"   
  error:
    cmds: # errors are reported in a special manner
      - "sleep f"
  silenterr:
    silent: true # errors are still reported in silent mode by design
    cmds:
      - "this-is-silent"
```


