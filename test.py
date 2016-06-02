#!/usr/bin/env python
#coding:utf-8

import time

class boardObject(object):
    def __init__(self, board = None):
        if board == None:
            row = [0, 0, 0, 0, 0, 0]
            self.board = list([list(row), list(row), list(row), list(row), list(row), list(row), list(row)])
            self.nbrPieces = list(row)
            self.player = 1
            self.turn = 1
        else:
            self.board = [list(board.board[0]), list(board.board[1]), list(board.board[2]), list(board.board[3]), list(board.board[4]), list(board.board[5]), list(board.board[6])]
            self.nbrPieces = list(board.nbrPieces)
            self.player = board.player
            self.turn = board.turn


    def printBoard(self):
        for y in reversed(range(0, 6)):
            for x in range(0, 7):
                if self.board[x][y] == 0:
                    print '.',
                elif self.board[x][y] == 1:
                    print 'X',
                elif self.board[x][y] == 2:
                    print '0',
                else:
                    raise ValueError("Values on the board can only by 0, 1 or 2")
            print
        print


    def addPiece(self, numRow):
        if numRow > 7 or numRow < 0:
            print "Value of row out of range"
            return False

        if self.nbrPieces[numRow] >= 6:
            print "Too many pieces on the row."
            return False

        if self.player == 1:
            self.board[numRow][self.nbrPieces[numRow]] = 1
        else:
            self.board[numRow][self.nbrPieces[numRow]] = 2

        self.nbrPieces[numRow] = self.nbrPieces[numRow] + 1
        self.player = 1 - self.player
        self.turn = self.turn + 1

        return True

    def hasWon(self, numRow):
        if numRow > 7 or numRow < 0:
            raise ValueError("Value of row out of range")


class connectFourObject(object):
    def __init__(self):
        self.boards = []
        self.nbrBoards = 0

    def addBoard(self, board = None):
        board = boardObject(board)
        self.boards.append(board)
        return board

    def removeBoard(self, board):
        self.boards.remove(board)

    def printBoards(self):
        for board in self.boards:
            for y in reversed(range(0, 6)):
                for x in range(0, 7):
                    if board.board[x][y] == 0:
                        print '.',
                    elif board.board[x][y] == 1:
                        print 'X',
                    elif board.board[x][y] == 2:
                        print '0',
                    else:
                        raise ValueError("Values on the board can only by 0, 1 or 2")
                print
            print
            print


game1 = connectFourObject()
board1 = game1.addBoard()

board1.printBoard()
board1.addPiece(1)

board1.printBoard()
board1.addPiece(1)

board1.printBoard()
board1.addPiece(1)

board1.printBoard()
board1.addPiece(1)
board2 = game1.addBoard(board1)

board1.printBoard()
board1.addPiece(1)

board1.printBoard()
board1.addPiece(1)

board1.printBoard()
board1.addPiece(0)

board1.printBoard()
board2.printBoard()

print board1.nbrPieces
print board2.nbrPieces

game1.printBoards()
game1.removeBoard(board2)

game1.printBoards()




