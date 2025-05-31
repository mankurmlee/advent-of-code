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
    distanceTo(other) {
        return Math.abs(other.x - this.x) + Math.abs(other.y - this.y)
    }
    toString() {
        return `Vec(${this.x},${this.y})`
    }
}

let me = {
    face: 0,
    pos: new Vec(0, 0),
    been: new Set(),
    twice: null,

    parse: function(instruction) {
        var dir = instruction.charAt(0)
        this.face += dir == 'L' ? 3 : 1
        this.face %= 4
        this.move((+instruction.slice(1)))
    },

    move: function(steps) {
        for (var i = 0; i < steps; i++) {
            if (this.face == 0) {
                this.pos.y--
            } else if (this.face == 1) {
                this.pos.x++
            } else if (this.face == 2) {
                this.pos.y++
            } else if (this.face == 3) {
                this.pos.x--
            }
            if (this.twice == null) {
                var posStr = this.pos.toString()
                if (this.been.has(posStr)) {
                    var v = new Vec(0, 0)
                    this.twice = this.pos.distanceTo(v)
                } else {
                    this.been.add(posStr)
                }
            }
        }
    },

    follow: function(directions) {
        for (var s of directions) {
            this.parse(s)
        }
    }
}

const input = fileinput()[0].split(', ')
me.follow(input)
console.log("Part 1:", me.pos.distanceTo(new Vec(0, 0)))
console.log("Part 2:", me.twice)
