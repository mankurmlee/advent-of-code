const fs = require('fs')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

function parse(data) {
    return data.map(l => l
        .split(/\s+/)
        .filter(s => s !== "")
        .map(n => +n))
}

function isTriangle(triplet) {
    t = [...triplet]
    t.sort((a, b) => a - b)
    return t[0] + t[1] > t[2]
}

function transpose(m) {
    let out = []
    let h = m.length
    for (let c = 0; c < 3; c++) {
        for (let r = 0; r < h; r += 3) {
            out.push([m[r][c], m[r+1][c], m[r+2][c]])
        }
    }
    return out
}

const input = parse(fileinput())
console.log("Part 1:", input.filter(isTriangle).length)

const transposed = transpose(input)
console.log("Part 2:", transposed.filter(isTriangle).length)
