# Pete Projects #

To give the child a name, Pete simply stands for PErformance TEst. 

## Why ##
Learning new languages like Elixir alone is already exciting but make it a bit more exciting by trying new frameworks and compare their speed against each other.
Because I'm a PHP (mostly Drupal) guy and just discovered Elixir I want to take it as opportunity to learn functional programming paradigms but also see how Elixir's (Erlang) concurrency model speeds out our traditional languages such as PHP. Therefore I took the Elixir web framework Phoenix Framework and try to compare it in a similar structure and setup to Slim Framework (PHP) and Phalcon (PHP as C extension) and see how the performance differs.

## What ##
I choose to do a very simplistic gallery app, with no model and database backend. 
Only router, controller, view and templates are used.

### Requriements ###
Image gallery:       
- Headline    
- 4 images (filenames as array/list with copyright info)    
- 3 Templates/partials (base, gallery, image)    
- CSS    
- No model and database backend (only controller / view / templates)    

### Framworks ###
- [Phoenix Framework (Elixir)](http://www.phoenixframework.org/)    
- [SlimPHP (PHP)](http://www.slimframework.com/)    
- [Phalcon (PHP, as compiled C extension)](https://phalconphp.com/en)    


## How ##
I'm not familar with any of the frameworks I have chosen here, so bare with me if some structure is a complete mess :)

See subdirectories for framework specific information.

## Test Results ##
I made multiple test runs on Raspberry Pi 2 (4 cores) and some Digitalocean VMs 
with 2 and 12 cores to check how the frameworks deal with multiple cores and concurrency.

**Update 2016-01-05:**     
After changing php.ini settings (see configuration settings below) PHP frameworks where 
able to improve by 10 - 40%. Updated test results. Would be great if we could figure out why Phoenix
is not using all system ressources on 12-core and it get's back to an interesting first place race with Phalcon at the top.

**Update 2016-01-06:**    
Ok, the history has to be rewritten. 
Because I used a [slow way](https://github.com/phoenixframework/phoenix/issues/1451) to 
include the images in the template file the performance was really bad. We tracked that issue down 
and now we really see the performance of Phoenix (Elixir/Erlang). I expected Phoenix to be faster
than PHP frameworks but that much of lead, I did not expect. Compared to the last test Phoenix speedup was 4x to 14x!
But take a look below yourself, unbelievable fast. And the craziest thing is, no matter how much load you throw at it, it seems 
to never return a bad request (3xx or 5xx) status code. (At least bandwidth seems to be maxed out so unable to throw more load at it :))

### Testing using wrk ###
Using [wrk](https://github.com/wg/wrk) as benchmarking tool with this or similar command:    
```
$ wrk -t4 -c100 -d60s --timeout 2000 http://ip-or-host.tld/gallery
```

### Language versions used for frameworks ###
PHP 5: PHP 5.6.14-0-deb8u1 as PHP-FPM, OPcache enabled, Nginx 1.6.2    
PHP 7: PHP 7.0.1 (self compiled) as PHP-FPM, OPcache enabled, Nginx 1.6.2    
Phoenix: Erlang 18.2.1 (SMP + hipe enabled), Elixir 1.2.0, Phoenix 1.1.0       

### Testing summary: Raspberry Pi 2 (Model B) ###
One nice thing about the Raspberry Pi is that the hardware is cheap and easy to 
get and test results can (hopefully) be reproduced and compared easier.           

| Framework      | Throughput (req/s) | Latency avg (ms) |     Stdev (ms) |
| :------------- | -----------------: | ---------------: | -------------: |
| Phoenix        |           2247.92  |           47.92  |         36.19  |
| Phalcon        |            586.41  |          170.22  |         18.68  |
| Slim (PHP 5.6) |            132.75  |          748.62  |        104.27  |
| Slim (PHP 7.0) |             27.71  |        >3500.00  |        553.34  |

Detailed results and specs: [results--raspberry-pi2.md](results--raspberry-pi2.md)

### Testing summary: Virtual Machine with 2 cores ###
Basic digitalocean VM with SSD and 2 cores, 2gb ram          

Pretty interesting that with fewer cores Phoenix is close but does not take the 
crown this time. 

| Framework      | Throughput (req/s) | Latency avg (ms) |     Stdev (ms) |
| :------------- | -----------------: | ---------------: | -------------: |
| Phalcon        |           1716.16  |           58.28  |         16.74  |
| Slim (PHP 7.0) |           1394.63  |           71.26  |         13.56  |
| Phoenix        |           1347.11  |           54.47  |         16.32  |
| Slim (PHP 5.6) |            693.11  |          143.40  |         20.78  |

Detailed results and specs: [results--2-core-vm.md](results--2-core-vm.md)

But now time to get really into multi core business and see who can handle concurrency best:

### Testing summary: Virtual Machine with 12 cores ###
App VM: Digitalocean VM with SSD and 12 cores, 32gb RAM     
Benchmark VM: Digitalocean VM with SSD and 4 cores, 8 gb RAM
Connected over 1 Gb/s

We do the same number of connections but with more threads first, then let's try
something insane and quadruple the number of connections and see what happens.

```
# ./wrk -t12 -c100 -d60s --timeout 1000   
```
| Framework      | Throughput (req/s) | Latency avg (ms) |     Stdev (ms) |
| :------------- | -----------------: | ---------------: | -------------: |
| Phoenix        |          38399.63  |            2.56  |          2.28  |
| Phalcon        |           9669.83  |           11.25  |         14.32  |
| Slim (PHP 7.0) |           7162.05  |           13.44  |          3.70  |
| Slim (PHP 5.6) |           3419.78  |           28.11  |          6.80  |

   
Ok, get a bit extreme and throw 12 threads with 400 connections onto the app server:   
```
# ./wrk -t12 -c400 -d180s --timeout 1000   
```
| Framework      | Throughput (req/s) | Latency avg (ms) |     Stdev (ms) |
| :------------- | -----------------: | ---------------: | -------------: |
| Phoenix        |          **  |          **  |         **  |
| Phalcon        |                 -  |               -  |             -  |
| Slim (PHP 7.0) |                 -  |               -  |             -  |
| Slim (PHP 5.6) |                 -  |               -  |             -  |

** 2016-01-06: will need to run this test sometime again, old value was 3000 req/s but it handles now way more. See insane test below :)      

Ok, so I really want to make Phoenix throw some errors now and do something insane. 48 treads 500 connections each.    
```
# ./wrk -t48 -c500 -d60s --timeout 1000   
```   
| Framework      | Throughput (req/s) | Latency avg (ms) |     Stdev (ms) |
| :------------- | -----------------: | ---------------: | -------------: |
| Phoenix        |          41162.79  |          43.55  |         111.75  |
| Phalcon        |                 -  |               -  |             -  |
| Slim (PHP 7.0) |                 -  |               -  |             -  |
| Slim (PHP 5.6) |                 -  |               -  |             -  |

This is just outstanding, take a look at the details below and look at the response times etc.

Detailed results, specs and comments: [results--12-core-vm.md](results--12-core-vm.md)


## Configuration ##
Some configs for PHP and Nginx can be found in ```configs``` directory:    

**Raspberry PI + 2 core VM configs for Nginx + PHP-FPM pool:**    
[configs/raspberry-pi2/nginx.conf](configs/raspberry-pi2/nginx.conf)      
[configs/raspberry-pi2/php-fpm_pool.d_www.conf](configs/raspberry-pi2/php-fpm_pool.d_www.conf)     

**12 Core VM**    
[configs/big-machine/nginx.conf](configs/big-machine/nginx.conf)      
[configs/big-machine/php-fpm_pool.d_www.conf](configs/big-machine/php-fpm_pool.d_www.conf)    

**Common configs (Nginx vhosts + php.ini used an all testsystems):**     
[configs/common/phalcon_nginx_vhost](configs/common/phalcon_nginx_vhost)        
[configs/common/slim_nginx_vhost](configs/common/slim_nginx_vhost)   
[configs/common/php.ini](configs/common/php.ini)    
Adding the following php.ini opcache-options improved above mentioned throughput for the PHP frameworks 
about 10 - 20 % of improvement (+ ~40% for Slim on 12-core!). Thanks to [@andresgutierrez](https://github.com/andresgutierrez) for suggesting these tweaks.
```
opcache.revalidate_freq = 300
opcache.enable_file_override = On
```
(Keep in mind that this revalidates opcache code only every 5 minutes, so you definetely want that option on 
production when you want to squeeze everything out and code does not change)   
[configs/common/phpinfo.html](configs/common/phpinfo.html)   

## Credits / inspiration  
I was inspired by these great guys but wanted to do it my way and see the results with a slightly more complex testapp.    
http://blog.onfido.com/using-cpus-elixir-on-raspberry-pi2/    
(title is a bit misleading because also PHP with FPM uses all CPUs pretty much, but more to that in a blog post)   
 
Comparison from Chris McCord (creator of Phoenix) with Rails:     
http://www.littlelines.com/blog/2014/07/08/elixir-vs-ruby-showdown-phoenix-vs-rails/    

Followup to above mentioned comparison but extended to Go, Ruby, NodeJS frameworks:   
https://github.com/mroth/phoenix-showdown      


License: [MIT](https://opensource.org/licenses/MIT)