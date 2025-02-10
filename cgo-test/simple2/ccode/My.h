#ifndef _LKH_MY
#define _LKH_MY

#include "LKH.h"
#include <limits.h>
#include <float.h>
#include <stdbool.h>
#include <stdint.h>

typedef enum  {
    CS_ALPHA,
    CS_DELAUNAY,
    CS_DELAUNAY_PURE,
    CS_NN,
    CS_POPMUSIC,
    CS_QUADRANT
} CandidateSetTypeEnum;

typedef enum {
    SUBPROBLEM_UNDEFINED,
    SUBPROBLEM_DELAUNAY,
    SUBPROBLEM_KARP,
    SUBPROBLEM_K_CENTER,
    SUBPROBLEM_K_MEANS,
    SUBPROBLEM_MOORE,
    SUBPROBLEM_ROHE,
    SUBPROBLEM_SIERPINSKI
} SubproblemSpecialEnum;

typedef enum {
    SPECIALSUBPROBLEM_UNDEFINED,
    SPECIALSUBPROBLEM_BORDERS,
    SPECIALSUBPROBLEM_COMPRESSED
} SubproblemSpecial2Enum;

typedef struct MyParameters {
    int AscentCandidates;
    uint32_t BackboneTrials;
    int Backtracking;
    int BWTSP_B;
    int BWTSP_Q;
    int BWTSP_L;
    enum CandidateSetTypes CandidateSetType;
    bool DelaunayPure;
    int MTSPDepot;
    double DistanceLimit;
    double Excess;
    int ExternalSalesmen;
    int ExtraCandidates;
    bool ExtraCandidateSetSymmetric;
    enum CandidateSetTypes ExtraCandidateSetType;
    bool Gain23Used;
    bool GainCriterionUsed;
    int InitialPeriod;
    int InitialStepSize;
    enum InitialTourAlgorithms InitialTourAlgorithm;
    double InitialTourFraction;
    int KickType;
    int Kicks;
    bool TSPTW_Makespan;
    int MaxBreadth;
    int MaxCandidates;
    bool CandidateSetSymmetric;
    int MaxSwaps;
    int MaxTrials;
    int MoveType;
    int MoveTypeSpecial;
    int MTSPMaxSize;
    int MTSPMinSize;
    enum Objectives MTSPObjective;
    int NonsequentialMoveType;
    double Optimum;
    int PatchingA;
    bool PatchingAExtended;
    bool PatchingARestricted;
    int PatchingC;
    bool PatchingCExtended;
    bool PatchingCRestricted;
    bool POPMUSIC_InitialTour;
    int POPMUSIC_MaxNeighbors;
    int POPMUSIC_SampleSize;
    int POPMUSIC_Solutions;
    int POPMUSIC_Trials;
    int MaxPopulationSize;
    int Precision;
    int Probability;
    enum RecombinationTypes Recombination;
    bool RestrictedSearch;
    int Runs;
    int Salesmen;
    int Scale;
    unsigned int Seed;
    bool Special;
    bool StopAtOptimum;
    bool Subgradient;
    int SubproblemSize;
    SubproblemSpecialEnum SubproblemSpecial;
    SubproblemSpecial2Enum SubproblemSpecial2;
    int SubsequentMoveType;
    bool SubsequentMoveTypeSpecial;
    bool SubsequentPatching;
    double TimeLimit;
    double TotalTimeLimit;
    int TraceLevel;
} MyParameters;

typedef struct {
    int Id;
    double X;
    double Y;
    double Z;
} NodeCoord;

typedef struct {
    uint32_t Capacity;
    uint32_t DemandDimension; // > 0
    uint32_t MTSPDepot; // > 0
    enum Types ProblemType;
    uint32_t Dimension;
    enum EdgeWeightTypes EdgeWeightType;
    NodeCoord *nodeCoords;
} MyProblem;

MyParameters createDefaultMyParameters();

// wrappers from ReadProblem.c due to static functions declaration
void CheckSpecificationPartWrapper(void);
void CreateNodesWrapper(void);
void Convert2FullMatrixWrapper(void);

#endif