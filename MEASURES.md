# Performance of various servers:
For performance testing, 10000 requests for the top news was sent using the http client, response time for  each request was obtained and mean, max and standard deviation measures were calculated.

1. Serial and No Caching
MEAN: 054.811309000000016ms
MAX: 1.041273s
STD: 035.542732631038916ms

2. Parallel and No Caching
MEAN: 057.83934499999994ms
MAX: 1.023297s
STD: 044.809368323175225ms

3. Serial and Caching
MEAN: 613.6530000000007us
MAX: 008.223ms
STD: 580.8058923521693us

4. Parallel and Caching 
MEAN: 664.3560000000009us
MAX: 009.913ms
STD: 656.4798072020185us

