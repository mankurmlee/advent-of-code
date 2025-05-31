const fs = require('fs')

function fileinput() {
    const data = fs.readFileSync(process.argv[2], 'utf8')
    return data.split('\n').slice(0, -1)
}

class Vec {
    constructor(x, y) {
        this.x = x
        this.y = y
    }
    clone() {
        return new Vec(this.x, this.y)
    }
    toString() {
        return `${this.x},${this.y}`
    }
}

class Keypad {
    constructor(dial) {
        this.dial = dial
        this.pos = new Vec(0, 0)
    }

    step(dir) {
        var v = this.pos.clone()
        if (dir == "U") {
            v.y--
        } else if (dir == "D") {
            v.y++
        } else if (dir == "L") {
            v.x--
        } else if (dir == "R") {
            v.x++
        }
        if (v.toString() in this.dial) {
            this.pos = v
        }
    }

    nextButton(dirs) {
        for (var d of dirs) {
            this.step(d)
        }
        return this.dial[this.pos.toString()]
    }

    getCode(input) {
        var nums = []
        for (var l of input) {
            var b = k.nextButton(l)
            nums.push(b)
        }
        return nums.join("")
    }
}

function square() {
    var dial = {}, i = 1
    for (var y = -1; y < 2; y++) {
        for (var x = -1; x < 2; x++) {
            var v = new Vec(x, y)
            dial[v.toString()] = i.toString()
            i++
        }
    }
    return dial
}

function diamond() {
    var dial = {}, i = 1
    for (var y = -2; y < 3; y++) {
        var xStart = Math.abs(y)
        var xEnd = 5 - xStart
        for (var x = xStart; x < xEnd; x++) {
            var v = new Vec(x, y)
            dial[v.toString()] = i.toString(16).toUpperCase()
            i++
        }
    }
    return dial
}

const input = fileinput()
let k = new Keypad(square())
console.log("Part 1:", k.getCode(input))
k = new Keypad(diamond())
console.log("Part 2:", k.getCode(input))
