# Pete Projects #

To give the child a name, Pete simply stands for PErformance TEst. 

## Why ##
Learning new languages like Elixir alone is already exciting but make it a bit more exciting by trying new frameworks and compare their speed against each other.
Because I'm a PHP (mostly Drupal) guy and just discovered Elixir I want to take it as opportunity to learn functional programming paradigms but also see how Elixir's (Erlang) concurrency model speeds out our traditional languages such as PHP. Therefore I took the Elixir web framework Phoenix Framework and try to compare it in a similar structure and setup to Slim Framework (PHP) and Phalcon (PHP as C extension) and see how the performance differs.

## What ##
I choose to do a very simplistic gallery app (no, this has nothing to do because I play with my raspberry pi camera these holidays).

### Requriements ###
Image gallery:       
- Headline    
- 4 images (filenames as array/list with copyright info)    
- 3 Templates/partials (base, gallery, image)    
- CSS    
- No model and database backend (only controller / view / templates)    

### Framworks ###
- Phoenix Framework (Elixir)    
- SlimPHP (PHP)    
- Phalcon (PHP, as compiled C extension)    


## How ##
I'm not familar with any of the frameworks I have chosen here, so bare with me if some structure is a complete mess :)

See subdirectories for specific information.

Todo:
For really deploying it and run it on production systems I also want to somehow deploy it with Docker or similar to get it installed fast on different VMs on DigitalOcean and/or AWS.

## Test Results ##
One nice thing about the Raspberry Pi is that the hardware is pretty the same and results can be reproduced.    
Hardware: Raspberry Pi 2 (4 cores, 1 GB Ram) Raspbian Jessie    
Network connection through 1Gbs RasPi <-> Router <-> USB dongle Macbook Pro   

### Testing using wrk ###
Using [wrk](https://github.com/wg/wrk) as benchmarking tool with this command:    
```
$ wrk -t4 -c100 -d60s --timeout 2000 http://ip-or-host/gallery
```

### Testing summary ###
| Framework      | Throughput (req/s) | Latency avg (ms) |     Stdev (ms) |
| :------------- | -----------------: | ---------------: | -------------: |
| Phoenix        |            575.23  |          173.46  |         12.52  |
| Phalcon        |            521.19  |          191.40  |         17.77  |
| Slim (PHP 5.6) |            117.72  |          844.07  |        101.49  |
| Slim (PHP 7.0) |             27.71  |        >3500.00  |        553.34  |

### Detailed results ###
#### Phoenix Framework (Elixir) ####
Erlang OTP 18, Elixir 1.2.0, Phoenix 1.1.0   
```
$ wrk -t4 -c100 -d60s --timeout 2000 http://pete-phoenix.pi:4001/gallery
Running 1m test @ http://pete-phoenix.pi:4001/gallery
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   173.46ms   12.52ms 261.63ms   71.21%
    Req/Sec   144.51     28.08   230.00     70.79%
  34569 requests in 1.00m, 65.08MB read
Requests/sec:    575.23
Transfer/sec:      1.08MB
```
Looks like a really good result right? Let's see how a not compiled framework like Slim Framework handles the task.

#### Slim Framework (PHP) ####
PHP 5.6.14-0+deb8u1 as FPM, Nginx 1.6.2
```
$ wrk -t4 -c100 -d60s --timeout 2000 http://pete-slim.pi/gallery
Running 1m test @ http://pete-slim.pi/gallery
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   844.07ms  101.49ms   1.25s    57.28%
    Req/Sec    31.25     18.32   118.00     66.15%
  7072 requests in 1.00m, 12.30MB read
Requests/sec:    117.72
Transfer/sec:    209.70KB
```
Yep, pretty poor, nearly 5 times slower than Elixir. Ok, PHP 7 should be twice as 
fast, so I compiled it on the Raspberry Pi but the results where even poorer (maybe I missed something?):    

PHP 7.0.1 as FPM, Nginx 1.6.2 
```
$ wrk -t4 -c100 -d60s --timeout 2000 http://pete-slim.pi/gallery
Running 1m test @ http://pete-slim.pi/gallery
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.50s   553.34ms   4.22s    93.33%
    Req/Sec     9.26      6.17    40.00     67.55%
  1665 requests in 1.00m, 2.94MB read
Requests/sec:     27.71
Transfer/sec:     50.04KB
```
Holy s..., that's really a bad performance for PHP 7. I expected 1.5 to 2x improvement 
towards PHP 5.6. I also wanted to try HHVM to compare but it seems not compile on ARM (Raspberry Pi).

Ok, not expecting much from Phalcon but what followed was really a big suprise, see below.

#### Phalcon Framework (a PHP C extension) ####
PHP 5.6.14-0+deb8u1 as FPM, Nginx 1.6.2
```
$ wrk -t4 -c100 -d60s --timeout 2000 http://pete-phalcon.pi/gallery
Running 1m test @ http://pete-phalcon.pi/gallery
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   191.42ms   17.77ms 415.41ms   82.85%
    Req/Sec   130.87     28.42   220.00     62.89%
  31322 requests in 1.00m, 55.14MB read
Requests/sec:    521.19
Transfer/sec:      0.92MB
```
This is really amazing, never thought that any PHP framework (as they are not compiled)
will come near Phoenix, but Phalcon is even faster as it is partly compiled.


## Credits / Inspiration ##  
I was inspired by these great guys but wanted to do it my way and see the results with a slightly more complex testapp.    
http://blog.onfido.com/using-cpus-elixir-on-raspberry-pi2/    
(title is a bit misleading because also PHP with FPM uses all CPUs pretty much, but more to that in a blog post)   
 
Comparison from Chris McCord (creator of Phoenix) with Ruby:     
http://www.littlelines.com/blog/2014/07/08/elixir-vs-ruby-showdown-phoenix-vs-rails/    

Followup to above mentioned comparison but extended to Go, Ruby, NodeJS frameworks:   
https://github.com/mroth/phoenix-showdown      
