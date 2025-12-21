#!/bin/bash
# Demo data setup script for THICC

echo "Resetting database..."
echo "yes" | ./thicc reset

echo ""
echo "Configuring settings..."
echo -e "lbs\nin\n70\n145" | ./thicc add 160 2024-01-01

echo ""
echo "Adding demo weight entries..."
# January - starting weight
./thicc add 160.5 2024-01-01
./thicc add 159.8 2024-01-05
./thicc add 159.2 2024-01-10
./thicc add 158.5 2024-01-15
./thicc add 158.0 2024-01-20
./thicc add 157.3 2024-01-25

# February - steady progress
./thicc add 156.8 2024-02-01
./thicc add 156.1 2024-02-05
./thicc add 155.5 2024-02-10
./thicc add 154.9 2024-02-15
./thicc add 154.2 2024-02-20
./thicc add 153.7 2024-02-25

# March - plateau
./thicc add 153.4 2024-03-01
./thicc add 153.2 2024-03-05
./thicc add 153.6 2024-03-10
./thicc add 153.0 2024-03-15
./thicc add 153.4 2024-03-20
./thicc add 152.8 2024-03-25

# April - breakthrough
./thicc add 152.3 2024-04-01
./thicc add 151.6 2024-04-05
./thicc add 151.0 2024-04-10
./thicc add 150.4 2024-04-15
./thicc add 149.8 2024-04-20
./thicc add 149.2 2024-04-25

# May - goal reached
./thicc add 148.7 2024-05-01
./thicc add 148.2 2024-05-05
./thicc add 147.5 2024-05-10

echo ""
echo "Demo data setup complete!"
echo ""
echo "Running thicc to show the table..."
echo ""
./thicc
