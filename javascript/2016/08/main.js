const fs = require('fs')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

function render() {
    let i = 0, j = WIDTH
    for (let y = 0; y < HEIGHT; y++) {
       console.log(SCREEN.slice(i, j).join(''))
       i = j
       j += WIDTH
    }
}

function rect(w, h) {
    for (let y = 0; y < h; y++) {
        const e = y*WIDTH + w
        for (let i = y*WIDTH; i < e; i++) {
            SCREEN[i] = '#'
        }
    }
}

function rotateRow(y, r) {
    r %= WIDTH
    if (r == 0) return
    let i = y * WIDTH
    const row = SCREEN.slice(i, i+WIDTH)
    i = WIDTH - r
    const rotated = row.slice(i).concat(row.slice(0, i))
    i = y * WIDTH
    for (let x = 0; x < WIDTH; x++) {
        SCREEN[i] = rotated[x]
        i++
    }
}

function rotateCol(x, r) {
    r %= HEIGHT
    if (r == 0) return
    const col = Array.from({length: HEIGHT}, (_, y) => SCREEN[y*WIDTH+x])
    let i = HEIGHT - r
    const rotated = col.slice(i).concat(col.slice(0, i))
    i = x
    for (let y = 0; y < HEIGHT; y++) {
        SCREEN[i] = rotated[y]
        i += WIDTH
    }
}

function load() {
    for (const l of fileinput()) {
        const words = l.split(" ")
        const [a, b] = l.match(/\d+/g).map(n => +n)
        if (words.length == 2) {
            rect(a, b)
        } else if (words[1] == "column") {
            rotateCol(a, b)
        } else if (words[1] == "row") {
            rotateRow(a, b)
        }
    }
}

const WIDTH = 50
const HEIGHT = 6
const SCREEN = Array.from({length: WIDTH*HEIGHT}, () => ' ')

load()
console.log("Part 1:", SCREEN.filter(e => e == '#').length)
console.log("Part 2:")
render()
