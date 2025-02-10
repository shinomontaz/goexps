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

extern void _Reset1();
extern void _Reset2();
extern void _Reset3();
extern void _Reset4();
extern void _Reset5();
extern void _Reset6();
extern void _Reset7();
extern void _Reset8();

bool ReadParametersFromStruct(MyParameters *params) 
{
    _Reset1();
    _Reset2();
    _Reset3();
    _Reset4();
    _Reset5();
    _Reset6();
    _Reset7();
    _Reset8();

    if (params->AscentCandidates < 2) {
        return false;
    }

    AscentCandidates = params->AscentCandidates;
    BackboneTrials = params->BackboneTrials;
    Backtracking = params->Backtracking;
    BWTSP_B = params->BWTSP_B;
    BWTSP_Q = params->BWTSP_Q;
    BWTSP_L = params->BWTSP_L;
    CandidateSetSymmetric = params->CandidateSetSymmetric;
    CandidateSetType = params->CandidateSetType;
    Crossover = ERXT;
    DelaunayPartitioning = 0;
    DelaunayPure = 0;
    DemandDimension = 1;
    DistanceLimit = params->DistanceLimit;
    Excess = params->Excess;
    ExternalSalesmen = params->ExternalSalesmen;
    ExtraCandidates = params->ExtraCandidates;
    ExtraCandidateSetSymmetric = params->ExtraCandidateSetSymmetric;
    ExtraCandidateSetType = params->ExtraCandidateSetType;
    Gain23Used = params->Gain23Used;
    GainCriterionUsed = params->GainCriterionUsed;
    GridSize = 1000000.0;
    InitialPeriod = params->InitialPeriod;
    InitialStepSize = params->InitialStepSize;
    InitialTourAlgorithm = params->InitialTourAlgorithm;
    InitialTourFraction = params->InitialTourFraction;
    KarpPartitioning = 0;
    KCenterPartitioning = 0;
    KMeansPartitioning = 0;
    Kicks = params->Kicks;
    KickType = params->KickType;
    MaxBreadth = params->MaxBreadth;
    MaxCandidates = params->MaxCandidates;
    MaxPopulationSize = params->MaxPopulationSize;
    MaxSwaps = params->MaxSwaps;
    MaxTrials = params->MaxTrials;
    MoorePartitioning = 0;
    MoveType = params->MoveType;
    MoveTypeSpecial = params->MoveTypeSpecial;
    MTSPDepot = params->MTSPDepot;
    MTSPMinSize = params->MTSPMinSize;
    MTSPMaxSize = params->MTSPMaxSize;
    MTSPObjective = params->MTSPObjective;
    NonsequentialMoveType = params->NonsequentialMoveType;
    Optimum = params->Optimum;
    PatchingA = params->PatchingA;
    PatchingC = params->PatchingC;
    PatchingAExtended = params->PatchingAExtended;
    PatchingARestricted = params->PatchingARestricted;
    PatchingCExtended = params->PatchingCExtended;
    PatchingCRestricted = params->PatchingCRestricted;
    Precision = params->Precision;
    POPMUSIC_InitialTour = params->POPMUSIC_InitialTour;
    POPMUSIC_MaxNeighbors = params->POPMUSIC_MaxNeighbors;
    POPMUSIC_SampleSize = params->POPMUSIC_SampleSize;
    POPMUSIC_Solutions = params->POPMUSIC_Solutions;
    POPMUSIC_Trials = params->POPMUSIC_Trials;
    Recombination = params->Recombination;
    RestrictedSearch = params->RestrictedSearch;
    RohePartitioning = 0;
    Runs = params->Runs; // >= 0
    Salesmen = params->Salesmen;
    Scale = params->Scale;
    Seed = params->Seed;
    SierpinskiPartitioning = 0;
    StopAtOptimum = params->StopAtOptimum;
    Subgradient = params->Subgradient;
    SubproblemBorders = 0;
    SubproblemsCompressed = 0;
    SubproblemSize = params->SubproblemSize;
    SubsequentMoveType = params->SubsequentMoveType;
    SubsequentMoveTypeSpecial = params->SubsequentMoveTypeSpecial;
    SubsequentPatching = params->SubsequentPatching;
    TimeLimit = params->TimeLimit;
    TotalTimeLimit = params->TotalTimeLimit;
    TraceLevel = params->TraceLevel;
    TSPTW_Makespan = params->TSPTW_Makespan;

    // Handle special cases
    if (params->Special) {
        Gain23Used = 0;
        KickType = 4;
        MaxSwaps = 0;
        MoveType = 5;
        MoveTypeSpecial = 1;
        MaxPopulationSize = 10;
    }

    // Handle subproblem special cases
    if (params->SubproblemSpecial == SUBPROBLEM_DELAUNAY) {
        DelaunayPartitioning = 1;
    } else if (params->SubproblemSpecial == SUBPROBLEM_KARP) {
        KarpPartitioning = 1;
    } else if (params->SubproblemSpecial == SUBPROBLEM_K_CENTER) {
        KCenterPartitioning = 1;
    } else if (params->SubproblemSpecial == SUBPROBLEM_K_MEANS) {
        KMeansPartitioning = 1;
    } else if (params->SubproblemSpecial == SUBPROBLEM_MOORE) {
        MoorePartitioning = 1;
    } else if (params->SubproblemSpecial == SUBPROBLEM_ROHE) {
        RohePartitioning = 1;
    } else if (params->SubproblemSpecial == SUBPROBLEM_SIERPINSKI) {
        SierpinskiPartitioning = 1;
    }

    // Handle subproblem special2 cases
    if (params->SubproblemSpecial2 == SPECIALSUBPROBLEM_BORDERS) {
        SubproblemBorders = 1;
    } else if (params->SubproblemSpecial2 == SPECIALSUBPROBLEM_COMPRESSED) {
        SubproblemsCompressed = 1;
    }

    // if (!ProblemFileName)
    //     eprintf("Problem file name is missing");

    // if (SubproblemSize == 0 && SubproblemTourFileName != 0)
    //     eprintf("SUBPROBLEM_SIZE specification is missing");

    // if (SubproblemSize > 0 && SubproblemTourFileName == 0)
    //     eprintf("SUBPROBLEM_TOUR_FILE specification is missing");

    // if (SubproblemSize > 0 && Salesmen > 1)
    //     eprintf("SUBPROBLEM specification not possible for SALESMEN > 1");

    if (CandidateSetType != DELAUNAY) 
        DelaunayPure = 0;

    if (CandidateSetType == DELAUNAY && params->DelaunayPure)
        DelaunayPure = 1;

    return true;
}
