const fs = require('fs')
const crypto = require('crypto')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

const doorId = fileinput()[0]
let one = []
let two = {}
let i = 0
while (true) {
    const salted = doorId + i
    const hash = crypto.createHash('md5').update(salted).digest('hex');
    if (hash.startsWith("00000")) {
        console.log(salted, hash)
        const k = hash[5]
        if (one.length < 8) {
            one.push(k)
        }
        if (k.match(/[0-7]/) && !two.hasOwnProperty(k)) {
            two[k] = hash[6]
            if (Object.keys(two).length == 8) {
                break
            }
        }
    }
    i++
}
console.log("Part 1:", one.join(""))
console.log("Part 2:", one.map((_, i) => two[i]).join(""))