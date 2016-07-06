package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "sync"
)

var wg sync.WaitGroup

const CONNECTX                  =   4 //nbr of disc to connect to win
const NBRROWS                   =   7
const NBRVALINCELL              =   3 // -1; 0; 1
const NBRLINES                  =   6
const NBRCELLS                  =   NBRROWS * NBRLINES
const NBRHIDDENLAYERS           =   20
const NBRTOTALLAYERS            =   NBRHIDDENLAYERS + 2 // +2 because of input and output
const NBRTOTALWEIGHTS           =   NBRTOTALLAYERS - 1
const NBRNEURONSPERHIDDENLAYER  =   NBRCELLS*5
const INPUT                     =   0
const OUPUT                     =   NBRTOTALLAYERS - 1
const LAMBDA                    =   0
var TABNBRNEURONSPERLAYER []int =   initNbrNeuronsPerLayer()

type CONNECTFOUR struct{
    moves[NBRCELLS] int
    board[NBRCELLS] int
    nbrDiscs[NBRROWS] int
    player int
    turn int
    wonBy int
}

type NEURALNETWORK struct{
    layers  [][]float64
    weights [][]float64

    player int
    win int
}

func initNbrNeuronsPerLayer() []int {
    tab   := make([]int, NBRTOTALLAYERS)
    tab[INPUT] = NBRCELLS * NBRVALINCELL
    for i := 1; i < NBRTOTALLAYERS - 1; i++ {
        tab[i] = NBRNEURONSPERHIDDENLAYER
    }
    tab[OUPUT] = NBRROWS

    return tab
}

func randFloat() float64 {
    return rand.Float64()
}

func initLayer(size int) []float64 {
    layer   := make([]float64, size)
    
    for i := 0; i < size; i++ {
        layer[i] = 0
    }

    return layer
}

func initLayers() [][]float64 {    
    layers   := make([][]float64, NBRTOTALLAYERS)
    
    for numLayer := 0; numLayer < NBRTOTALLAYERS; numLayer++ { // 1 to MAX -1 because input and output don't have the same size compared to the others layers
        layers[numLayer] = initLayer(TABNBRNEURONSPERLAYER[numLayer]) 
    }

    return layers
}

func initWeight(size int) []float64 {
    weight   := make([]float64, size)
    
    for i := 0; i < size; i++ {
        weight[i] = randFloat()  * 2 - 1
    }
    
    return weight
}

func initWeights() [][]float64 {    
    weights   := make([][]float64, NBRTOTALWEIGHTS)
    
    for numWeight := 0; numWeight < NBRTOTALWEIGHTS; numWeight++ {
        weights[numWeight] = initWeight(TABNBRNEURONSPERLAYER[numWeight] * TABNBRNEURONSPERLAYER[numWeight + 1]) 
    }
    
    return weights
}

func initNeuralNetwork(pNeuralNetwork *NEURALNETWORK, player int) {
    pNeuralNetwork.win      =   0
    pNeuralNetwork.player   =   player
    
    pNeuralNetwork.layers   =   initLayers()
    pNeuralNetwork.weights  =   initWeights()
}

func calculValueNeurons(pNeuralNetwork *NEURALNETWORK, numLayer int, numNeuron int){
    pNeuralNetwork.layers[numLayer][numNeuron] = 0
    for numWeight := 0; numWeight < TABNBRNEURONSPERLAYER[numLayer - 1]; numWeight++ {
        pNeuralNetwork.layers[numLayer][numNeuron] += pNeuralNetwork.weights[numLayer - 1][TABNBRNEURONSPERLAYER[numLayer - 1] * numNeuron + numWeight] * pNeuralNetwork.layers[numLayer - 1][numWeight]
    }
    pNeuralNetwork.layers[numLayer][numNeuron] = 1 / (1 + math.Exp(pNeuralNetwork.layers[numLayer][numNeuron] * -1))
    if pNeuralNetwork.layers[numLayer][numNeuron] > 0.5 {
        pNeuralNetwork.layers[numLayer][numNeuron] = 1
    } else {
        pNeuralNetwork.layers[numLayer][numNeuron] = 0
    }

}

func calculValueLayers(pNeuralNetwork *NEURALNETWORK, numLayer int) {
    for numNeuron := 0; numNeuron < TABNBRNEURONSPERLAYER[numLayer]; numNeuron++ {
        calculValueNeurons(pNeuralNetwork, numLayer, numNeuron)
    }
}

func calculOutputNeuralNetwork(pNeuralNetwork *NEURALNETWORK, pConnectFour *CONNECTFOUR) {
    for numCell := 0; numCell < NBRCELLS * NBRVALINCELL; numCell += NBRVALINCELL {
        if pNeuralNetwork.player == pConnectFour.board[numCell / NBRVALINCELL] {
            pNeuralNetwork.layers[INPUT][numCell    ] = 1
            pNeuralNetwork.layers[INPUT][numCell + 1] = 0
            pNeuralNetwork.layers[INPUT][numCell + 2] = 0
        }else if pConnectFour.board[numCell / NBRVALINCELL] != 0 {
            pNeuralNetwork.layers[INPUT][numCell    ] = 0
            pNeuralNetwork.layers[INPUT][numCell + 1] = 1
            pNeuralNetwork.layers[INPUT][numCell + 2] = 0
        } else {
            pNeuralNetwork.layers[INPUT][numCell    ] = 0
            pNeuralNetwork.layers[INPUT][numCell + 1] = 0
            pNeuralNetwork.layers[INPUT][numCell + 2] = 1
        }
        //~ pNeuralNetwork.Layer[INPUT][numCell] = pConnectFour.board[numCell]
    }
    for numLayer := 1; numLayer < NBRTOTALLAYERS; numLayer++ {
        calculValueLayers(pNeuralNetwork, numLayer)
    }
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
    for x := 0; x < NBRCELLS; x++ {
        pConnectFour.moves[x] = -1
    }
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
    if checkWinRow(pConnectFour, numRow) || checkWinDiagonals(pConnectFour, numRow) || checkWinLine(pConnectFour, pConnectFour.nbrDiscs[numRow] - 1){
        pConnectFour.wonBy = pConnectFour.player
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
    var total float64 = 0
    var subTotal float64 = 0
    for layer := 0; layer < NBRTOTALLAYERS; layer++ {
        fmt.Println(pNeuralNetwork.layers[layer])
    }
    for row := 0; row < NBRROWS; row++ {
        if pConnectFour.nbrDiscs[row] >= NBRLINES {
            continue
        }
        pNeuralNetwork.layers[OUPUT][row] += LAMBDA
        total += pNeuralNetwork.layers[OUPUT][row]
    }
    valueRand := randFloat() * total
    for row := 0; row < NBRROWS; row++ {
        if pConnectFour.nbrDiscs[row] >= NBRLINES {
            continue
        }
        subTotal += pNeuralNetwork.layers[OUPUT][row]
        if subTotal >= valueRand {
            return row
        }
    }
    
    fmt.Println()
    fmt.Println()
    fmt.Println()
    fmt.Println(pNeuralNetwork.player)
    fmt.Println(pNeuralNetwork.layers[INPUT])
    fmt.Println(pNeuralNetwork.layers[OUPUT])
    fmt.Println(total)
    fmt.Println(valueRand)
    fmt.Println(subTotal)
    printBoard(pConnectFour)
    return NBRROWS
}


func playAGameHvB(pNeuralNetwork *NEURALNETWORK, humanPlayer int){
    
    var pConnectFour *CONNECTFOUR
    var connectFour CONNECTFOUR
    pConnectFour = &connectFour
    var row int
    
    initBoard(pConnectFour)
    for pConnectFour.turn <= NBRCELLS {
        printBoard(pConnectFour)
        if pConnectFour.player == humanPlayer{
            for {
                fmt.Scanf("%d", &row)
                if makeMove(pConnectFour, row) {
                    break
                }
            }
        } else {
            calculOutputNeuralNetwork(pNeuralNetwork, pConnectFour)
            row = chooseMove(pNeuralNetwork, pConnectFour)
            makeMove(pConnectFour, row)
        }
        if isWin(pConnectFour, row){
            printBoard(pConnectFour)
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
    for pConnectFour.turn < NBRCELLS {
        //~ printBoard(pConnectFour)
        if pConnectFour.player == pNeuralNetwork1.player{
            //~ calculOutput(pNeuralNetwork1, pConnectFour)
            //~ row = chooseMove(pNeuralNetwork1, pConnectFour)
        } else {
            //~ calculOutput(pNeuralNetwork2, pConnectFour)
            //~ row = chooseMove(pNeuralNetwork2, pConnectFour)
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
        //~ calculOutput(pNeuralNetwork, pConnectFour)
        //~ row = chooseMove(pNeuralNetwork, pConnectFour)
        makeMove(pConnectFour, row)
        if isWin(pConnectFour, row){
            //~ fmt.Println("Win")
            //~ fmt.Println(pConnectFour.wonBy)
            break
        }
        nextPlayer(pConnectFour)
        pNeuralNetwork.player *= -1
    }
    if pConnectFour.turn >= NBRCELLS {
        //~ printBoard(pConnectFour)
    } else {
        //~ trainNeuralNetworkBvHim(pNeuralNetwork, pConnectFour)
    }
    //~ fmt.Println(pConnectFour.moves)
}


func main(){
    rand.Seed(time.Now().UTC().UnixNano())
    
    var pNeuralNetwork1 *NEURALNETWORK
    var neuralNetwork1 NEURALNETWORK
    pNeuralNetwork1 = &neuralNetwork1
    
    var pNeuralNetwork2 *NEURALNETWORK
    var neuralNetwork2 NEURALNETWORK
    pNeuralNetwork2 = &neuralNetwork2
    
    initNeuralNetwork(pNeuralNetwork1, 1)
    initNeuralNetwork(pNeuralNetwork2, -1)
    //~ calculOutputNeuralNetwork(pNeuralNetwork1)
    //~ fmt.Println(neuralNetwork1.layers[OUPUT])
    start := time.Now()
    for cpt := 0;cpt < 1; cpt++ {
        playAGameHvB(pNeuralNetwork1, -1)
        //~ playAGameBvB(pNeuralNetwork1, pNeuralNetwork2)
        //~ fmt.Println(cpt)
        //~ playAGameBvHim(pNeuralNetwork1)
        //~ fmt.Println(cpt)
    }
    elapsed := time.Since(start)
    fmt.Printf("%s\n", elapsed)
    
}
