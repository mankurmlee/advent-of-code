const fs = require('fs')

function fileinput() {
    return fs.readFileSync(process.argv[2], 'utf8')
        .split('\n').slice(0, -1)
}

function setupLinks() {
    const input = fileinput()
        .filter(s => s.match(/(bot|output) \d+/g))
        .map(s => s.match(/(bot|output) \d+/g))
        .filter(n => n.length == 3)
    for (let o of input) {
        LINKS[o[0]] = o.slice(1)
    }
}

function setupBots() {
    const input = fileinput()
        .map(s => s.match(/\d+/g))
        .filter(n => n.length == 2)
        .slice(1)
    for (let [v, k] of input) {
        botPush("bot " + k, v)
    }
}

function botPush(k, v) {
    if (BINS[k]) {
        BINS[k].push(v)
    } else {
        BINS[k] = [v]
    }
    return BINS[k].length
}

function run() {
    const [targetLo, targetHi] = ordered(fileinput()[0].split(' '))
    const q = Object.entries(BINS)
        .filter(e => e[1].length == 2)
        .map(e => e[0])
    while (q.length > 0) {
        let b = q.shift()
        let [lo, hi] = ordered(BINS[b])
        delete BINS[b]
        if (lo == targetLo && hi == targetHi) {
            console.log("Part 1:", b.split(' ')[1])
        }
        dests = LINKS[b]
        k = dests[0]
        if (botPush(k, lo) > 1 && k.startsWith("bot")) {
            q.push(k)
        }
        k = dests[1]
        if (botPush(k, hi) > 1 && k.startsWith("bot")) {
            q.push(k)
        }
    }
}

const ordered = ([a, b]) => +a < +b ? [a, b] : [b, a]

const LINKS = {}
setupLinks()

const BINS = {}
setupBots()

run()
const prod = +BINS['output 0'] * +BINS['output 1']* +BINS['output 2']
console.log("Part 2:", prod)
