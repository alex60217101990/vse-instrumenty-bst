package helpers

import "time"

func Int32(v int32) *int32 {
    return &v
}

func Int32Ptr(v *int32) int32 {
    if v != nil {
        return *v
    }
    return 0
}

func Int64(v int64) *int64 {
    return &v
}

func Int64Ptr(v *int64) int64 {
    if v != nil {
        return *v
    }
    return 0
}

func Uint64(v uint64) *uint64 {
    return &v
}

func Uint64Ptr(v *uint64) uint64 {
    if v != nil {
        return *v
    }
    return 0
}

func Duration(v time.Duration) *time.Duration {
    return &v
}

func DurationPtr(v *time.Duration) time.Duration {
    if v != nil {
        return *v
    }
    return 0
}

func Bool(v bool) *bool {
    return &v
}

func BoolPtr(v *bool) bool {
    if v != nil {
        return *v
    }
    return false
}

func Int(v int) *int {
    return &v
}

func IntPtr(v *int) int {
    if v != nil {
        return *v
    }
    return 0
}

func TimeToTimePtr(t time.Time) *time.Time {
    return &t
}

func TimePtrToTime(t *time.Time) (emptyTime time.Time) {
    if t != nil {
        return *t
    }
    return emptyTime
}

func String(v string) *string {
    return &v
}

func StringPtr(v *string) string {
    if v != nil {
        return *v
    }
    return ""
}
