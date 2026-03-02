add1 n = n+1

fib 0 = 0
fib 1 = 1
fib n = fib (n-1) + fib (n-2)

mr l = reverse l


|| American call option via binomial tree (Turner-style, correct comprehension)

option_value spot duration strike sigma =
  v
    where
    || model parameters (per step)
    u      = exp sigma
    d      = exp (-sigma)

    r_step = 0.00016
    disc   = exp (-r_step)
    p      = (exp r_step - d) / (u - d)
    q      = 1 - p

    || stock price and exercise payoff
    stock t i = spot * u^i * d^(t - i)

    exercise t i
      = e, if e > 0
      = 0, otherwise
        where e = stock t i - strike

    terminal =
      [ exercise duration i | i <- [0..duration] ]

    || discounted hold value
    hold vd vu = disc * (p * vu + q * vd)

    || one backward time step (zip first, then match tuples)
    step (x:xs) t =
      [ max2 exi (hold vd vu)
      | ((vd,vu), exi)
          <- zip2 (zip2 (x:xs) xs)
                   [ exercise t i | i <- [0..t] ]
      ]

    || times to fold over
    ts
      = [], if duration = 0
      = [duration - 1, duration - 2 .. 0], otherwise

    result = foldl step terminal ts
    (v:ignore) = result


main =
  option_value 184.94 31 170 0.04

mms