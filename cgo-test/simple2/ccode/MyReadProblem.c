#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <limits.h>
#include <float.h>
#include "LKH.h"
#include "GainType.h"
#include "Genetic.h"
#include "My.h"

// Function to read node coordinates from a structure
static void MyRead_NODE_COORD_SECTION(NodeCoord *nodeCoords, int numNodes) {
    Node *N;
    int Id, i;

    CheckSpecificationPartWrapper();

    if (!FirstNode) {
        CreateNodesWrapper();
    }
    N = FirstNode;

    do {
        N->V = 0;
    } while ((N = N->Suc) != FirstNode);

    if (ProblemType == HPP)
        Dimension--;

    for (i = 0; i < numNodes; i++) {
        Id = nodeCoords[i].Id;

        if (Id <= 0 || Id > Dimension) {
            eprintf("NODE_COORD_SECTION: Node number out of range: %d", Id);
        }

        N = &NodeSet[Id];
        if (N->V == 1) {
            eprintf("NODE_COORD_SECTION: Node number occurs twice: %d", N->Id);
        }

        N->V = 1;
        N->X = nodeCoords[i].X;
        N->Y = nodeCoords[i].Y;

        if (CoordType == THREED_COORDS) {
            N->Z = nodeCoords[i].Z;
        }
    }

    N = FirstNode;
    do
        if (!N->V && N->Id <= Dim)
            break;
    while ((N = N->Suc) != FirstNode);
    if (!N->V)
        eprintf("NODE_COORD_SECTION: No coordinates given for node %d",
                N->Id);
    if (ProblemType == HPP)
        Dimension++;
    if (Asymmetric)
        Convert2FullMatrixWrapper();
}