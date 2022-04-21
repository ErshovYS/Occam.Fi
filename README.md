# Occam.Fi  

The solution is to receive data from different streams and combine them into one stream (or into a queue). As soon as the next minute starts, the last price value is displayed.  
  
In general, to draw beautiful bars, you need the maximum, minimum, first and last values, but since only one value is specified in the example, it would be more correct to take the last one.