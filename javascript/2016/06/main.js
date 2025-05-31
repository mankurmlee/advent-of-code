const fs = require('fs')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

function decode(fn, i) {
    let counts = {}
    for (const l of lines) {
        const ch = l[i]
        if (counts[ch]) {
            counts[ch]++
        } else {
            counts[ch] = 1
        }
    }
    return Object.entries(counts).reduce(fn)[0]
}

const lines = fileinput()
const mostFreq = (a, b) => a[1] > b[1] ? a : b
const leastFreq = (a, b) => a[1] < b[1] ? a : b

console.log("Part 1:", lines[0].split("")
    .map((_, i) => decode(mostFreq, i)).join(""))

console.log("Part 2:", lines[0].split("")
    .map((_, i) => decode(leastFreq, i)).join(""))
