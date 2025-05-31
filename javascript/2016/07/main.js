const fs = require('fs')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

function abba(s) {
    const n = s.length - 3
    for (let i = 0; i < n; i++) {
        if (s[i] == s[i+3] && s[i+1] == s[i+2] && s[i] != s[i+1]) {
            return true
        }
    }
    return false
}

function aba(s) {
    let found = []
    const n = s.length - 2
    for (let i = 0; i < n; i++) {
        if (s[i] == s[i+2] && s[i] != s[i+1]) {
            found.push(s[i] + s[i+1])
        }
    }
    return found
}

function tls(s) {
    const seqs = s.replace(/\[|\]/g, " ").split(" ")
    const n = seqs.length
    let odd = false
    for (let i = 0; i < n; i += 2) {
        if (abba(seqs[i])) {
            odd = true
            break
        }
    }
    let even = false
    for (let i = 1; i < n; i += 2) {
        if (abba(seqs[i])) {
            even = true
            break
        }
    }
    return odd && !even
}

function ssl(s) {
    const seqs = s.replace(/\[|\]/g, " ").split(" ")
    const n = seqs.length
    let abas = new Set()
    for (let i = 0; i < n; i += 2) {
        for (const k of aba(seqs[i])) {
            abas.add(k)
        }
    }
    for (let i = 1; i < n; i += 2) {
        for (const k of aba(seqs[i])) {
            if (abas.has(k[1] + k[0])) {
                return true
            }
        }
    }
    return false
}

console.log("Part 1:", fileinput().filter(ip => tls(ip)).length)
console.log("Part 2:", fileinput().filter(ip => ssl(ip)).length)