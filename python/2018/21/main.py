results = set()
target = 0
while True:
    mod = target | 65536
    target = 1397714
    while True:
        target += mod & 255
        target &= 16777215
        target *= 65899
        target &= 16777215
        if 256 > mod:
            break
        i = 0
        while (i + 1) * 256 <= mod:
            i += 1
        mod = i
    # First target is part 1
    # Last unique target is part 2
    if target not in results:
        results.add(target)
        print(target)
