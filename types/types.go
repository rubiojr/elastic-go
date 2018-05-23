package types

type ClusterHealth struct {
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
	ClusterName                 string  `json:"cluster_name"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	RelocatingShards            int     `json:"relocating_shards"`
	Status                      string  `json:"status"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	TimedOut                    bool    `json:"timed_out"`
	UnassignedShards            int     `json:"unassigned_shards"`
}
