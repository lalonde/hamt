#include "textflag.h"

        TEXT ·popcnt(SB),$0-16
        NOP
        NOP
        MOVQ x+0(FP), BX
        POPCNTQ BX, BX
        MOVQ BX, ret+8(FP)
        RET
