const fs = require('fs')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

function checksum(s) {
    let counts = {}
    for (const ch of s) {
        if (!(ch in counts)) {
            counts[ch] = 0
        }
        counts[ch]++
    }
    let countList = Object.entries(counts)
    countList.sort((a, b) => (
        a[1] != b[1] ? b[1] - a[1] : asc(a[0]) - asc(b[0])
    ))
    return countList.slice(0, 5).map(c => c[0]).join("")
}

function isRealRoom(s) {
    const letters = s.replace(/[^a-z]/g, "")
    return checksum(letters.slice(0, -5)) == letters.slice(-5)
}

function rotateLetter(ch, n) {
    if (ch == ' ') {
        return ch
    }
    const a = asc('a')
    return String.fromCharCode((asc(ch) - a + n) % 26 + a)
}

function rotateName(s) {
    const enc = s.replace(/-\d+.*$/, '').replace(/-/g, ' ')
    const rot = sectorID(s) % 26
    return enc.split("")
        .map((c) => rotateLetter(c, rot))
        .join("")
}

const asc = (ch) => ch.charCodeAt(0)
const sectorID = (s) => +s.match(/\d+/)[0]

const realRooms = fileinput().filter(isRealRoom)
console.log("Part 1:", realRooms.map(sectorID).reduce((a, b) => a + b))

const room = realRooms.filter(
    (s) => rotateName(s) === "northpole object storage"
)[0]
console.log("Part 2:", sectorID(room))
