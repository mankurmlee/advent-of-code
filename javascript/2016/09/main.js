const fs = require('fs')

function fileinput() {
    return fs.readFileSync(process.argv[2], 'utf8')
        .split('\n').slice(0, -1)
}

function len(s, recursive=false) {
    let l = 0, i = 0, marker, c
    const n = s.length
    while (i < n) {
        if (s[i] !== '(') {
            l++
            i++
        } else {
            marker = s.substring(i).match(/^\((\d+)x(\d+)\)/)
            c = +marker[1]
            i += marker[0].length + c
            if (recursive) {
                c = len(s.substring(i-c, i), true)
            }
            l += c * +marker[2]
        }
    }
    return l
}

console.log("Part 1:", len(fileinput()[0]))
console.log("Part 2:", len(fileinput()[0], true))