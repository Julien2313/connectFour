package main

import "fmt"
//~ import (
        //~ "fmt/Printf" 
        //~ Printf "printf"
       //~ )
//~ import "fmt/Printf"
import "math"
import "math/rand"
import "time"

const CONNECTX                  = 4 //nbr of isc to connect to win
const NBRROWS                   = 7
const NBRLINES                  = 6
const NBRCELLS                  = NBRROWS * NBRLINES
const NBRVALUEINCELL            = 3 // 1, -1 or 0
const NBRHIDDENNEURONSPERINPUT  = 4
const LAMBDA                    = 0

type CONNECTFOUR struct{
    moves[NBRCELLS] int
    board[NBRCELLS] int
    nbrDiscs[NBRROWS] int
    player int
    turn int
}

type NEURALNETWORK struct{
    input[NBRCELLS][NBRVALUEINCELL] float32
    coefInputToHiddenNeurons[NBRCELLS][NBRVALUEINCELL][NBRHIDDENNEURONSPERINPUT] float32
    valueHiddenNeurons[NBRCELLS][NBRHIDDENNEURONSPERINPUT] float32
    coefHiddenNeuronsToOutput[NBRCELLS][NBRHIDDENNEURONSPERINPUT][NBRROWS] float32
    output[NBRROWS] float32
    player int
    win int
}

func initNeuralNetwork(pNeuralNetwork *NEURALNETWORK, player int){
    var cell, numValueInCell, numHiddenNeurons, row int
    pNeuralNetwork.player = player
    pNeuralNetwork.win = 0

    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    
    var maxValueFloat32 float32 = math.MaxFloat32 / 10000000000000000000000000000000000000 / NBRCELLS
    
    for cell = 0; cell < NBRCELLS; cell++ {
        for numValueInCell = 0; numValueInCell < NBRVALUEINCELL; numValueInCell++ {
            for numHiddenNeurons = 0; numHiddenNeurons < NBRHIDDENNEURONSPERINPUT; numHiddenNeurons++ {
                pNeuralNetwork.coefInputToHiddenNeurons[cell][numValueInCell][numHiddenNeurons] = r1.Float32() * maxValueFloat32
            }
        }
    }

    for cell = 0; cell < NBRCELLS; cell++ {
        for numHiddenNeurons = 0; numHiddenNeurons < NBRHIDDENNEURONSPERINPUT; numHiddenNeurons++ {
            for row = 0; row < NBRROWS; row++ {
                pNeuralNetwork.coefHiddenNeuronsToOutput[cell][numHiddenNeurons][row] = r1.Float32() * maxValueFloat32
            }
        }
    }
}

func initInput(pNeuralNetwork *NEURALNETWORK, pConnectFour *CONNECTFOUR){
    
    for cell := 0; cell < NBRCELLS; cell++ {
        if pNeuralNetwork.player == pConnectFour.board[cell] {
            pNeuralNetwork.input[cell][0] = 1
            pNeuralNetwork.input[cell][1] = 0
            pNeuralNetwork.input[cell][2] = 0
        }else if pConnectFour.board[cell] != 0 {
            pNeuralNetwork.input[cell][0] = 0
            pNeuralNetwork.input[cell][1] = 1
            pNeuralNetwork.input[cell][2] = 0
        } else {
            pNeuralNetwork.input[cell][0] = 0
            pNeuralNetwork.input[cell][1] = 0
            pNeuralNetwork.input[cell][2] = 1
        }
    }
}

func calculValueHiddenNeurons(pNeuralNetwork *NEURALNETWORK) {
    var cell, cell2, numValue, numHiddenNeurons int
    var coef, value float32
    for cell = 0; cell < NBRCELLS; cell++ {
        for numHiddenNeurons = 0; numHiddenNeurons < NBRHIDDENNEURONSPERINPUT; numHiddenNeurons++ {
            pNeuralNetwork.valueHiddenNeurons[cell][numHiddenNeurons] = 0
            for cell2 = 0; cell2 < NBRCELLS; cell2++ {
                for numValue = 0; numValue < NBRVALUEINCELL; numValue++ {
                    coef = pNeuralNetwork.coefInputToHiddenNeurons[cell2][numValue][numHiddenNeurons]
                    value = pNeuralNetwork.input[cell2][numValue]
                    pNeuralNetwork.valueHiddenNeurons[cell][numHiddenNeurons] += coef * value
                }
            }
            //~ fmt.Printf("%f\n", pNeuralNetwork.valueHiddenNeurons[cell][numHiddenNeurons])
        }
    }
}

func calculValueOuput(pNeuralNetwork *NEURALNETWORK, pConnectFour *CONNECTFOUR) {
    var row, cell, numHiddenNeurons int
    var coef, value float32
    for row = 0; row < NBRROWS; row++ {
        pNeuralNetwork.output[row] = 0
        if pConnectFour.nbrDiscs[row] >= NBRLINES {
            continue
        }
        for cell = 0; cell < NBRCELLS; cell++ {
            for numHiddenNeurons = 0; numHiddenNeurons < NBRHIDDENNEURONSPERINPUT; numHiddenNeurons++ {
                coef = pNeuralNetwork.coefHiddenNeuronsToOutput[cell][numHiddenNeurons][row]
                value = pNeuralNetwork.valueHiddenNeurons[cell][numHiddenNeurons]
                pNeuralNetwork.output[row] += value * coef
            }
        }
        //~ fmt.Printf("%f\n", pNeuralNetwork.output[row])
    }
}

func getOutput(pNeuralNetwork *NEURALNETWORK, pConnectFour *CONNECTFOUR){
    initInput(pNeuralNetwork, pConnectFour)
    calculValueHiddenNeurons(pNeuralNetwork)
    calculValueOuput(pNeuralNetwork, pConnectFour)
    //~ fmt.Println(pNeuralNetwork.output)
    //~ fmt.Printf("%f, %f\n", pNeuralNetwork.output[0], pNeuralNetwork.output[1])

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

func nextPlayer(pConnectFour *CONNECTFOUR){
    pConnectFour.player *= -1
}

func initBoard(pConnectFour *CONNECTFOUR){
    pConnectFour.player = 1
    pConnectFour.turn = 0
}

func canAddDisc(pConnectFour *CONNECTFOUR, numRow int) bool{
    if pConnectFour.nbrDiscs[numRow] < NBRLINES {
        return true
    }
    return false
}

func addDisc(pConnectFour *CONNECTFOUR, numRow int){
    pConnectFour.board[pConnectFour.nbrDiscs[numRow] * NBRROWS + numRow] = pConnectFour.player
    pConnectFour.nbrDiscs[numRow]++
    pConnectFour.moves[pConnectFour.turn] = numRow
    pConnectFour.turn++
}

func checkWinLine(pConnectFour *CONNECTFOUR, numLine int) bool{
    var combo int = 0
    for x := 0; x < NBRROWS; x++ {
        if pConnectFour.board[numLine * NBRROWS + x] == pConnectFour.player {
            combo++
            if combo == CONNECTX {
                return true
            }
        } else {
            combo = 0
        }
    }
    return false
}
func checkWinRow(pConnectFour *CONNECTFOUR, numRow int) bool{
    var combo int = 0
    for y := 0; y < NBRLINES; y++ {
        if pConnectFour.board[y * NBRROWS + numRow] == pConnectFour.player {
            combo++
            if combo == CONNECTX {
                return true
            }
        } else {
            combo = 0
        }
    }
    return false
}
func checkWinDiagonals(pConnectFour *CONNECTFOUR, numRow int) bool{
    var combo  int = 0
    var height int = pConnectFour.nbrDiscs[numRow] - 1

    var x int = max(0, numRow - height)
    var y int = max(0, height - numRow)
    for y < NBRLINES && x < NBRROWS {
        if pConnectFour.board[y * NBRROWS + x] == pConnectFour.player {
            combo++
            if combo == CONNECTX {
                return true
            }
        } else {
            combo = 0
        }
        y++
        x++ 
    }

    x = max(0, max(0, height + numRow) - (NBRLINES - 1))
    y = min((NBRLINES - 1), max(0, height + numRow))

    for y >= 0 && x < NBRROWS {
        if pConnectFour.board[y * NBRROWS + x] == pConnectFour.player {
            combo++
            if combo == CONNECTX {
                return true
            }
        } else {
            combo = 0
        }
        y--
        x++ 
    }
    return false
}


func isWin(pConnectFour *CONNECTFOUR, numRow int) bool{
    if checkWinRow(pConnectFour, numRow) {
        return true
    }

    if checkWinLine(pConnectFour, pConnectFour.nbrDiscs[numRow] - 1) {
        return true
    }
    
    if checkWinDiagonals(pConnectFour, numRow) {
        return true
    }
        
    return false
}

func makeMove(pConnectFour *CONNECTFOUR, numRow int) bool{
    if canAddDisc(pConnectFour, numRow) {
        addDisc(pConnectFour, numRow)
        return true        
    }
    return false
}

func printBoard(pConnectFour *CONNECTFOUR){
    var pos, disc int
    for y := NBRLINES - 1; y >= 0; y-- {
        for x := 0; x < NBRROWS; x++ {
            pos = y * NBRROWS + x
            disc = pConnectFour.board[pos]
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

func chooseMove(pNeuralNetwork *NEURALNETWORK, pConnectFour *CONNECTFOUR) int {
    var total float32 = 0
    var subTotal float32 = 0
    for row := 0; row < NBRROWS; row++ {
        if pConnectFour.nbrDiscs[row] >= NBRLINES {
            continue
        }
        pNeuralNetwork.output[row] += LAMBDA
        total += pNeuralNetwork.output[row]
    }
    seed := rand.NewSource(time.Now().UnixNano())
    rand := rand.New(seed)
    valueRand := rand.Float32() * total
    for row := 0; row < NBRROWS; row++ {
        if pConnectFour.nbrDiscs[row] >= NBRLINES {
            continue
        }
        subTotal += pNeuralNetwork.output[row]
        if subTotal > valueRand {
            return row
        }
    }
    
    return NBRROWS - 1
}

func playAGameHvB(pNeuralNetwork *NEURALNETWORK, humanPlayer int){
    
    var pConnectFour *CONNECTFOUR
    var connectFour CONNECTFOUR
    pConnectFour = &connectFour
    var row int
    
    initBoard(pConnectFour)
    for {
        printBoard(pConnectFour)
        if pConnectFour.player == humanPlayer{
            for {
                fmt.Scanf("%d", &row)
                if makeMove(pConnectFour, row) {
                    break
                }
            }
        } else {
            getOutput(pNeuralNetwork, pConnectFour)
            row = chooseMove(pNeuralNetwork, pConnectFour)
            makeMove(pConnectFour, row)
        }
        printBoard(pConnectFour)
        if isWin(pConnectFour, row){
            fmt.Println("Win")
            break
        }
        nextPlayer(pConnectFour)
    }
}

func playAGameBvB(pNeuralNetwork1 *NEURALNETWORK, pNeuralNetwork2 *NEURALNETWORK){
    
    var pConnectFour *CONNECTFOUR
    var connectFour CONNECTFOUR
    pConnectFour = &connectFour
    var row int
    
    initBoard(pConnectFour)
    for pConnectFour.turn <= NBRCELLS {
        //~ printBoard(pConnectFour)
        if pConnectFour.player == pNeuralNetwork1.player{
            getOutput(pNeuralNetwork1, pConnectFour)
            row = chooseMove(pNeuralNetwork1, pConnectFour)
        } else {
            getOutput(pNeuralNetwork2, pConnectFour)
            row = chooseMove(pNeuralNetwork2, pConnectFour)
        }
        makeMove(pConnectFour, row)
        if isWin(pConnectFour, row){
            //~ fmt.Println("Win")
            if pConnectFour.player == pNeuralNetwork1.player{
                pNeuralNetwork1.win++
            } else {
                pNeuralNetwork2.win++
            }
            break
        }
        nextPlayer(pConnectFour)
    }
    if pConnectFour.turn > NBRCELLS {
        printBoard(pConnectFour)
    }
}

func playAGameBvHim(pNeuralNetwork *NEURALNETWORK){
    
    var pConnectFour *CONNECTFOUR
    var connectFour CONNECTFOUR
    pConnectFour = &connectFour
    var row int
    
    initBoard(pConnectFour)
    for pConnectFour.turn < NBRCELLS {
        //~ printBoard(pConnectFour)
        getOutput(pNeuralNetwork, pConnectFour)
        row = chooseMove(pNeuralNetwork, pConnectFour)
        makeMove(pConnectFour, row)
        if isWin(pConnectFour, row){
            //~ fmt.Println("Win")
            pNeuralNetwork.win++
            break
        }
        nextPlayer(pConnectFour)
        pNeuralNetwork.player *= -1
    }
    if pConnectFour.turn >= NBRCELLS {
        printBoard(pConnectFour)
    }
}


func main(){
    //~ humanPlayer := 1
    //~ var pNeuralNetwork *NEURALNETWORK
    //~ var neuralNetwork NEURALNETWORK
    //~ pNeuralNetwork = &neuralNetwork
    //~ initNeuralNetwork(pNeuralNetwork, humanPlayer * -1)
    //~ playAGameHvB(pNeuralNetwork, humanPlayer)
    
    var pNeuralNetwork1 *NEURALNETWORK
    var neuralNetwork1 NEURALNETWORK
    pNeuralNetwork1 = &neuralNetwork1

    var pNeuralNetwork2 *NEURALNETWORK
    var neuralNetwork2 NEURALNETWORK
    pNeuralNetwork2 = &neuralNetwork2
    
    initNeuralNetwork(pNeuralNetwork1, 1)
    initNeuralNetwork(pNeuralNetwork2, -1)
    
    
    start := time.Now()
    for cpt := 0;cpt < 10000; cpt++ {
        //~ playAGameBvB(pNeuralNetwork1, pNeuralNetwork2)
        //~ fmt.Println(cpt)
        playAGameBvHim(pNeuralNetwork1)
    }
    elapsed := time.Since(start)
    fmt.Printf("%s\n", elapsed)
    fmt.Printf("%d\n", pNeuralNetwork1.win)
    fmt.Printf("%d\n", pNeuralNetwork2.win)
}
