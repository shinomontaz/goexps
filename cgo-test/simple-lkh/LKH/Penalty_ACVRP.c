#include "LKH.h"
#include "Segment.h"

GainType Penalty_ACVRP() {
    static Node *StartRoute = 0;
    Node *N, *NextN, *CurrentRoute;
    GainType DemandSum, DistanceSum, P = 0, P2 = 0;
    int Forward = SUCC(Depot)->Id != Depot->Id + DimensionSaved;

    if (!StartRoute)
        StartRoute = Depot;
    if (StartRoute->Id > DimensionSaved)
        StartRoute -= DimensionSaved;
    N = StartRoute;
    int cnt_loaded_cars = 0;

    do {
        int Size = 0;
        CurrentRoute = N;
        DemandSum = 0;
        do {
            if (N->Id <= Dim && N != Depot) {
                if ((DemandSum += N->Demand) > Capacity)
                    P += DemandSum - Capacity;
                if (P + P2 > CurrentPenalty ||
                    (P + P2 == CurrentPenalty && CurrentGain <= 0)) {
                    StartRoute = CurrentRoute;
                    return CurrentPenalty + (CurrentGain > 0);
                }
                Size++;
            }
            N = Forward ? SUCC(N) : PREDD(N);
        } while (N->DepotId == 0);
        if (MTSPMinSize >= 1 && Size < MTSPMinSize)
            P += MTSPMinSize - Size;
        if (Size > MTSPMaxSize)
            P += Size - MTSPMaxSize;
        if (DemandSum > 0) {
            cnt_loaded_cars += 1;
        }
        if (DemandSum > 0 && DemandSum < CapacityLowerBound) {
            P2 += CapacityLowerBound - DemandSum;
        }
    } while (N != StartRoute);
    P2 += Capacity * (Salesmen - cnt_loaded_cars);
    if (P < PenaltyLowerBound) {
        P = PenaltyLowerBound;
    }
    if (P2 < PenaltyLowerBoundCapacity) {
        P2 = PenaltyLowerBoundCapacity;
    }
//    printf("total Penalty for under constraints = %lld \n\n", P2);

    return P + P2;
}
