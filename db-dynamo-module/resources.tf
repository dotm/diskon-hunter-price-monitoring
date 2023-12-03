#Put resources here.
#If this files become to long, you can move related resources to their own files.

#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
}

#Note:
#Free tier: 25 GB of Storage, 25 provisioned Write Capacity Units (WCU), 25 provisioned Read Capacity Units (RCU)
#Calculation:
#1 RCU  = 2 eventually consistent reads of up to 4 KB/s. (5KB -> 1RCU)
#1 RCU  = 1 strongly consistent read of up to 4 KB/s. (5KB -> 2RCU)
#2 RCUs = 1 transactional read request (one read per second) for items up to 4 KB. (5KB -> 4RCU)
#1 WCU  = 1 standard write of up to 1 KB/s. (5KB -> 5WCU)
#2 WCUs = 1 transactional write request (one write per second) for items up to 1 KB. (5KB -> 10WCU)
