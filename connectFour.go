package main

import "fmt"

const CONNECTX  = 4
const NBRROWS   = 7
const NBRLINES  = 6
const NBRCELLS  = NBRROWS * NBRLINES

type CONNECTFOUR struct{
    board[NBRCELLS] int
    nbrDiscs[NBRROWS] int
    player int
}

func max(a, b int) int{
    if a > b{
        return a
    }
    return b
}

func min(a, b int) int{
    if a < b{
        return a
    }
    return b
}

func nextPlayer(connectFour *CONNECTFOUR){
    connectFour.player *= -1
}

func initBoard(connectFour *CONNECTFOUR){
    connectFour.player = 1    
}

func canAddDisc(connectFour *CONNECTFOUR, numRow int) bool{
    if connectFour.nbrDiscs[numRow] < NBRLINES {
        return true
    }
    return false
}

func addDisc(connectFour *CONNECTFOUR, numRow int){
    connectFour.board[connectFour.nbrDiscs[numRow] * NBRROWS + numRow] = connectFour.player
    connectFour.nbrDiscs[numRow]++
}

func checkWinLine(connectFour *CONNECTFOUR, numLine int) bool{
    combo := 0
    for x := 0; x < NBRROWS; x++ {
        if connectFour.board[numLine * NBRROWS + x] == connectFour.player {
            combo++
        } else {
            combo = 0
        }
        
        if combo == CONNECTX {
            return true
        }
    }
    return false
}
func checkWinRow(connectFour *CONNECTFOUR, numRow int) bool{
    combo := 0
    for y := 0; y < NBRLINES; y++ {
        if connectFour.board[y * NBRROWS + numRow] == connectFour.player {
            combo++
        } else {
            combo = 0
        }
        
        if combo == CONNECTX {
            return true
        }
    }
    return false
}

func checkWinDiagonals(connectFour *CONNECTFOUR, numRow int) bool{
    combo   := 0
    height  := connectFour.nbrDiscs[numRow] - 1

    x := max(0, numRow - height)
    y := max(0, height - numRow)
    for y < NBRLINES && x < NBRROWS {
        if connectFour.board[y * NBRROWS + x] == connectFour.player {
            combo++
        } else {
            combo = 0
        }
        if combo == CONNECTX {
            return true
        }
        y++
        x++ 
    }

    x = max(0, max(0, height + numRow) - (NBRLINES - 1))
    y = min((NBRLINES - 1), max(0, height + numRow))


    for y >= 0 && x < NBRROWS {
    //~ fmt.Println(x)
    //~ fmt.Println(y)
    //~ fmt.Println()
        if connectFour.board[y * NBRROWS + x] == connectFour.player {
            combo++
        } else {
            combo = 0
        }
        if combo == CONNECTX {
            return true
        }
        //~ fmt.Println(combo)
        y--
        x++ 
    }
    return false
}


func isWin(connectFour *CONNECTFOUR, numRow int) bool{
    if checkWinRow(connectFour, numRow) {
        return true
    }

    if checkWinLine(connectFour, connectFour.nbrDiscs[numRow] - 1) {
        return true
    }
    
    if checkWinDiagonals(connectFour, numRow) {
        return true
    }
        
    return false
}

func makeMove(connectFour *CONNECTFOUR, numRow int) bool{
    if canAddDisc(connectFour, numRow) {
        addDisc(connectFour, numRow)
        return true        
    }
    return false
}

func printBoard(connectFour *CONNECTFOUR){
    var pos, disc int
    for y := NBRLINES - 1; y >= 0; y-- {
        for x := 0; x < NBRROWS; x++ {
            pos = y * NBRROWS + x
            disc = connectFour.board[pos]
            if disc == 1 {
                fmt.Print("X ")
            } else if disc == -1{
                fmt.Print("0 ")
            } else {
                fmt.Print(". ")
            }
        }
        fmt.Println()
    }
    fmt.Println()
    fmt.Println()
}


func main(){

    var connectFour CONNECTFOUR
    var row int
    initBoard(&connectFour)
    printBoard(&connectFour)
    
    for {
        for {
            fmt.Scanf("%d", &row)
            if makeMove(&connectFour, row) {
                break
            }
        }
        printBoard(&connectFour)
        if isWin(&connectFour, row){
            fmt.Println("Win")
            break
        }
        nextPlayer(&connectFour)
    }
}
