Core switch ports + panel
=========================

port    | panel  | deviceport  | name
--------|--------|------------------
Gi1/0/1 |        |             | DEBUG
Gi1/0/2 |        |             | DEBUG/trunk
Gi1/0/3 |        | onboard eth | vin
Gi1/0/5 |        | eth0        | ap-cf-srv
Gi1/0/6 | 01.06A | eth0        | ap-cf-f-l
Gi1/0/7 | 01.10A | eth0        | ap-cf-f-r
Gi1/0/8 | 01.09A | eth0        | ap-cf-b1
Gi1/0/9 | 01.18A | eth0        | ap-cf-a1
Gi1/0/10| 01.08A | eth0        | ap-cf-a2
Gi1/0/11| 01.20A | eth0        | ap-cf-c1
Gi1/0/17| 01.08B | Gi1/0/8     | reception-sw
Gi1/0/18| unknown| Gi1/0/8     | team-sw
Gi1/0/19| 01.19A | Gi1/0/8     | vocsw-A
Gi1/0/20| 01.16A | Gi1/0/8     | vocsw-B
Gi1/0/21| 01.21A | Gi1/0/8     | vocsw-C
Gi1/0/22| unknown| Gi1/0/8     | vocsw-D
Gi1/0/24| 01.18B |             | presenter-A 
Gi1/0/25| 01.17A |             | presenter-B
Gi1/0/26| 01.20B |             | presenter-C
Gi1/0/27| 01.09A | eth0        | cambox-B
Gi1/0/28| 01.05A | eth0        | cambox-C
Gi1/0/48|techpark| unknown     | techpark switch
Gi1/0/49|unknown | gi0/49      | f0sw (MM fiber)
Gi1/0/50|unknown | n/a         | ipacct (SM fiber)
Te0/2   |        | enp1s0f4d1  | vin (MM fiber)

