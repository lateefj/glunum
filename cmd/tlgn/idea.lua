print("Idea of how to use glunum")
x = {8, 2, -9, 15, 4}
std = stat.StdDev(x, nil)
print("Standard Deviation is " .. std)
weights = {2, 2, 6, 7, 1}
wstd = stat.StdDev(x, weights)
print("Weighted Standard Deviation is " .. wstd)

mean = stat.Mean(x, nil)
print("Mean is " .. mean)
wm = stat.Mean(x, weights)
print("Weighted Mean " .. wm)

skew = stat.Skew(x, nil)
print("Skew is " .. mean)
ws = stat.Skew(x, weights)
print("Weighted Skew " .. ws)

variance = stat.Variance(x, nil)
print("Variance is " .. variance)
wv = stat.Variance(x, weights)
print("Weighted Variance " .. wv)
