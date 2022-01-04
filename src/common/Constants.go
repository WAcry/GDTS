package common

const (
	// JOB_SAVE_DIR job save dir
	JOB_SAVE_DIR = "/cron/jobs/"

	// JOB_KILLER_DIR job killer dir
	JOB_KILLER_DIR = "/cron/killer/"

	// JOB_LOCK_DIR job lock dir
	JOB_LOCK_DIR = "/cron/lock/"

	// JOB_WORKER_DIR service register dir
	JOB_WORKER_DIR = "/cron/workers/"

	// JOB_EVENT_SAVE save job event
	JOB_EVENT_SAVE = 1

	// JOB_EVENT_DELETE delete job event
	JOB_EVENT_DELETE = 2

	// JOB_EVENT_KILL kill job event
	JOB_EVENT_KILL = 3
)
