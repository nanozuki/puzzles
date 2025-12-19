# Nonogram

Algorithm:

- Find all possible patterns for each row and column based on clues
- Use all candidate patterns to find force filled and force empty cells in a
  line
- Propagate the filled and empty cells to other lines until no more cells can be
  propagated
- Use backtracking to try each branch if the propagation cannot fill all cells

## Benchmark

For current version of tests:

### Test P001

Nonogram with 5 rows and 5 columns, Solved Nonogram use 26 steps, 0.00s:

```
+-+-+-+-+-+
|· o o o ·|
|· o · o ·|
|· o o o ·|
|· · o · ·|
|· · o o ·|
+-+-+-+-+-+
```

### P060

Nonogram with 10 rows and 10 columns, Solved Nonogram use 108 steps, 0.00s:

```
+-+-+-+-+-+-+-+-+-+-+
|· · · o o o · · o o|
|o · o · o · o o · ·|
|· o o · o o · · · o|
|· · o · · o o · o ·|
|· · · o o · · o o o|
|o · o · o · · · · o|
|o o o o · o · · o ·|
|· · · · o · o o · ·|
|· · · o · · · o · ·|
|· o o · o · · o o ·|
+-+-+-+-+-+-+-+-+-+-+
```

### P097

Nonogram with 15 rows and 15 columns, Solved Nonogram use 170 steps, 0.00s:

```
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|· · · · · · · · · · · o · o ·|
|· · · · · · · o o o o · o o ·|
|· · · · · · o o o o o o · o ·|
|· · · · · · o o o · o · o · o|
|· · · · · · o o · · o o o o ·|
|· · · · · · o o · · o o o · ·|
|· · · · · o o o o o o · o · o|
|· · · · o o · o o · · · o o o|
|· · · o o · o · · · · o o o ·|
|· · o o · · o · · · o o · o ·|
|· o o · · · o o · · · · · o ·|
|o · o · · · · o · · · · · o ·|
|o · o · · · · o o · · · · o ·|
|· o o · · · · · o · · · · · o|
|· · · · · · · · o o o o o o o|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

# P136

New Nonogram with 15 rows and 20 columns, Solved Nonogram use 132 steps, 0.00s:

```
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|· · · o o o o · o o o · · o o o · · · ·|
|· · · · o o o o · o o o o o o o o o o o|
|· · o o · o o o · · o · o o · o o o o ·|
|· · o o o o o · o o o o · o o o · o o o|
|· · · · o o · o · o o · o o · · · o o ·|
|· · o o o o o · o · · o o · · o o · · o|
|o o o o o · o o · · o o o o o o · o o o|
|o · o · o o · o o o · o o o o o o o o ·|
|· · · · · o o o · · · · · o · o · · · ·|
|· · · · · o o o · · · · · · o o · · · ·|
|· · · · · o o · · · · · · · · o o · · ·|
|· · · · · o o · · · · · · · · o o · · ·|
|· · · · · o o o o · · · · · · o o · · ·|
|· · · · o o o o · o · · · · o o o · · ·|
|· · · · o o · o · o · · · o o · o · · ·|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

### P150

New Nonogram with 15 rows and 20 columns, Solved Nonogram use 234 steps,
0.00s:

```
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|· · · · · · · o · · o · · o · · o · o ·|
|· · · · o o o o o · o · o o o · o o o ·|
|· · · o · · · o o o o o · · o o · · o o|
|· · · o · o · · o · o o o o · · o · o ·|
|· · · · o o · · o o o o o · · · · o o o|
|· · o o o · o o · o · · o · o · · o o o|
|· o · · · o · · o · · o o o o · o · o o|
|· o · o · · o · o · · o · · · o o · o ·|
|· · o o · · o o o o o · · · · · o o o ·|
|· o o o o o · o · · · o · o · · o · o o|
|o · · · o · o · · o · · o o o · o o o o|
|o · o · · o o · o o · o o o · o o · o ·|
|· o o · · o · o o o o o o o o o o o · ·|
|· · o · o · o · · · o o o o o o · o o ·|
|· · · o o o · · o · o · o o o o · o · o|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```
