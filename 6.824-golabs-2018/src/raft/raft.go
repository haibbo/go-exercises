package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"labrpc"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// import "bytes"
// import "labgob"

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type LogEntry struct {
	Term int
	Cmd  interface{}
}

const (
	// Follower state is the initialization state
	Follower int8 = iota
	Candidate
	Leader
)

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	state           int8 // Three states, Follower, Candidate and Leader
	electionTimer   *time.Timer
	electionResetCh chan struct{}
	commitCh        chan struct{} // used to commit log
	applyCh         chan ApplyMsg
	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	// Persistent state on all servers
	currentTerm int        //latest term server has seen (initialized to 0 on first boot, increases monotonically)
	votedFor    int        //candidateId that received vote in current term (or null if none)
	log         []LogEntry //log entries; each entry contains command for state machine, and term when entry was received by leader (first index is 1)

	// Volatile state on all servers:
	commitIndex int //index of highest log entry known to be committed (initialized to 0, increases monotonically)
	lastApplied int // index of highest log entry applied to state machine (initialized to 0, increases monotonically)

	// Volatile state on leaders: (Reinitialized after election)
	nextIndex  []int // for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	matchIndex []int // for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	term = rf.currentTerm
	isleader = false
	if rf.state == Leader {
		isleader = true
	}
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
//  AppendEntries RPC arguments  stucture
//
type AppendEntriesArgs struct {
	Term         int        //leader’s term
	LeaderID     int        //so follower can redirect clients
	PrevLogIndex int        //index of log entry immediately preceding new ones
	PrevLogTerm  int        //term of prevLogIndex entry
	Entries      []LogEntry //log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        //leader’s commitIndex
}

//
//  AppendEntries RPC reply stucture
//
type AppendEntriesReply struct {
	Term    int  //currentTerm, for leader to update itself
	Success bool //true if follower contained entry matching PrevLogIndex and prevLogTerm
}

func (rf *Raft) GetLastLogTermIndex() (int, int) {
	lastIndex := len(rf.log)
	if lastIndex <= 1 {
		return -1, -1
	}
	return rf.log[lastIndex-1].Term, lastIndex - 1
}

//
// example AppendEntries RPC handler.
//
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	// Your code here (2A, 2B).
	// discard out of date AppendEntries
	DPrintf("[%d:%d] Recv AppendEntries with %d log\n", rf.me, rf.currentTerm, len(args.Entries))
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		DPrintf("[%d:%d] Drop out of date AppendEntries with Term %d ", rf.me, rf.currentTerm, args.Term)
		return
	}
	rf.electionResetCh <- struct{}{}
	// If RPC request or response contains term T > currentTerm: set currentTerm = T, convert to follower
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = Follower
		rf.votedFor = args.LeaderID
	}
	//Reply false if log doesn’t contain an entry at prevLogIndex whose term matches prevLogTerm
	if len(rf.log)-1 < args.PrevLogIndex {
		reply.Term = rf.currentTerm
		reply.Success = false
		DPrintf("[%d:%d] AppendEntries return false2", rf.me, rf.currentTerm)
		return
	}
	DPrintf("[%d:%d] Heartbeat %v %v", rf.me, rf.currentTerm, args, rf.log)
	// If an existing entry conflicts with a new one (same index but different terms), delete the existing entry and all that follow it
	if rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		DPrintf("[%d:%d] local log %d's term %d is not same request term %d", rf.me, rf.currentTerm, args.PrevLogIndex, rf.log[args.PrevLogIndex].Term, args.PrevLogTerm)
		rf.log = rf.log[:args.PrevLogIndex]
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}
	// overwrite logs after PreLogIndex
	rf.log = append(rf.log[:args.PrevLogIndex+1], args.Entries...)
	DPrintf("[%d:%d] LeaderCommit:%d rf.commitIndex:%d", rf.me, rf.currentTerm, args.LeaderCommit, rf.commitIndex)
	//If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	_, lastIndex := rf.GetLastLogTermIndex()
	if args.LeaderCommit > rf.commitIndex {
		if args.LeaderCommit < lastIndex {
			rf.commitIndex = args.LeaderCommit
		} else {
			rf.commitIndex = lastIndex
		}
		rf.commitCh <- struct{}{}
	}
	reply.Term = rf.currentTerm
	reply.Success = true
	return
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int //candidate’s term
	CandidateID  int // candidate requesting vote
	LastLogIndex int // index of candidate’s last log entry (§5.4)
	LastLogTerm  int // term of candidate’s last log entry (§5.4)
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int  //currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.VoteGranted = false
	DPrintf("[%d:%d] Recv RequestVote from %d Term: %d\n", rf.me, rf.currentTerm, args.CandidateID, args.Term)
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}
	// If RPC request or response contains term T > currentTerm: set currentTerm = T, convert to follower
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = Follower
		rf.votedFor = -1
	}
	// If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote
	if rf.votedFor == -1 || rf.votedFor == args.CandidateID {
		lastTerm, lastIndex := rf.GetLastLogTermIndex()
		DPrintf("[%d:%d] compare [%d:%d] [%d:%d] %v",
			rf.me, rf.currentTerm,
			args.LastLogTerm, args.LastLogIndex, lastTerm, lastIndex, rf.log)
		if args.LastLogTerm > lastTerm ||
			(args.LastLogTerm == lastTerm && args.LastLogIndex >= lastIndex) {
			rf.votedFor = args.CandidateID
			reply.Term = rf.currentTerm
			reply.VoteGranted = true
			DPrintf("[%d:%d] compare2 true", rf.me, rf.currentTerm)
		}
	}

	reply.Term = rf.currentTerm
	rf.electionResetCh <- struct{}{}
	DPrintf("[%d:%d] End RequestVote reply.VoteGranted %d\n", rf.me, rf.currentTerm, reply.VoteGranted)
	return
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	index := -1
	term := -1
	isLeader := true
	// Your code here (2B).
	if rf.state != Leader {
		return index, term, false
	}
	log := LogEntry{
		Term: rf.currentTerm,
		Cmd:  command,
	}
	rf.log = append(rf.log, log)
	rf.nextIndex[rf.me] = len(rf.log)
	rf.matchIndex[rf.me] = len(rf.log) - 1
	DPrintf("[%d:%d] Add new log %v, total len = %d", rf.me, rf.currentTerm, log, len(rf.log))
	return len(rf.log) - 1, log.Term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
	_, leader := rf.GetState()
	DPrintf("[%d-%d] Call Kill, Leader is %v\n", rf.me, rf.currentTerm, leader)
}

func (rf *Raft) UpdateCommitIndex() {
	matchIndex := make([]int, len(rf.matchIndex))
	copy(matchIndex, rf.matchIndex)
	sort.Ints(matchIndex)

	N := matchIndex[(len(rf.peers)+1)/2-1]
	DPrintf("[%d:%d] N: %d, old commitIndex =  %d", rf.me, rf.currentTerm, N, rf.commitIndex)
	//If there exists an N such that N > commitIndex, a majority of matchIndex[i] ≥ N,
	// and log[N].term == currentTerm:	set commitIndex = N
	if N > rf.commitIndex && rf.log[N].Term == rf.currentTerm {
		rf.commitIndex = N
		rf.commitCh <- struct{}{}
	}
	DPrintf("[%d:%d] N: %d, new  commitIndex = %d", rf.me, rf.currentTerm, N, rf.commitIndex)
}
func (rf *Raft) Heartbeat() {
	for {
		DPrintf("[%d:%d] Enter Heartbeat", rf.me, rf.currentTerm)
		if rf.state != Leader {
			return
		}
		if len(rf.electionTimer.C) == 0 {
			rf.electionResetCh <- struct{}{}
		}
		for i := 0; i < len(rf.peers); i++ {
			if i == rf.me {
				continue
			}
			go func(i int) {
				var reply AppendEntriesReply
				rf.mu.Lock()
				logs := make([]LogEntry, len(rf.log)-rf.nextIndex[i])
				copy(logs, rf.log[rf.nextIndex[i]:])

				DPrintf("[%d:%d] send AppendEntries to %d with %d log %v commitIndex:%d, %d %v",
					rf.me, rf.currentTerm, i, len(logs), logs, rf.commitIndex, rf.nextIndex[i], rf.log)
				Args := &AppendEntriesArgs{
					Term:         rf.currentTerm,
					LeaderID:     rf.me,
					PrevLogIndex: rf.nextIndex[i] - 1,
					PrevLogTerm:  rf.log[rf.nextIndex[i]-1].Term,
					Entries:      logs,
					LeaderCommit: rf.commitIndex,
				}
				rf.mu.Unlock()
				if rf.sendAppendEntries(i, Args, &reply) {
					rf.mu.Lock()
					if reply.Success {
						rf.nextIndex[i] = Args.PrevLogIndex + len(Args.Entries) + 1
						rf.matchIndex[i] = rf.nextIndex[i] - 1
						DPrintf("[%d:%d] increase %d nextIndex %d, len log: %d len Entries: %d", rf.me, rf.currentTerm, i, rf.nextIndex[i], len(rf.log), len(Args.Entries))
						rf.UpdateCommitIndex()
					} else {
						if reply.Term == Args.Term {
							rf.nextIndex[i]--
							DPrintf("Try new nextIndex %d for %d", rf.nextIndex[i], i)
						} else if reply.Term > Args.Term {
							rf.currentTerm = reply.Term
							rf.state = Follower
							rf.votedFor = -1
						}
					}
					rf.mu.Unlock()
				}
			}(i)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
func (rf *Raft) CommitLogs() {
	for {
		select {
		case <-rf.commitCh:
			rf.mu.Lock()
			for i := rf.lastApplied; i <= rf.commitIndex; i++ {
				msg := ApplyMsg{
					CommandValid: true,
					Command:      rf.log[i].Cmd,
					CommandIndex: i,
				}
				DPrintf("[%d:%d] Commit log %v", rf.me, rf.currentTerm, rf.log[i])
				rf.applyCh <- msg
				rf.lastApplied = i
			}
			rf.mu.Unlock()
		}
	}
}
func (rf *Raft) LeaderElection() {
	rf.state = Follower
	rf.electionTimer = time.NewTimer(time.Millisecond * time.Duration(rand.Int()%300+500))
	for {
		select {
		case <-rf.electionTimer.C:
			DPrintf("[%d:%d] election timeout, current state %d\n", rf.me, rf.currentTerm, rf.state)
			rf.votedFor = rf.me
			rf.currentTerm++
			rf.state = Candidate
			Args := &RequestVoteArgs{
				Term:         rf.currentTerm,
				CandidateID:  rf.me,
				LastLogIndex: len(rf.log) - 1,
				LastLogTerm:  rf.log[len(rf.log)-1].Term,
			}
			var votes int
			votes++
			for i := 0; i < len(rf.peers); i++ {
				if i == rf.me {
					continue
				}
				go func(i int) {
					var reply RequestVoteReply
					DPrintf("[%d:%d] Send RequestVote", rf.me, rf.currentTerm)
					if rf.sendRequestVote(i, Args, &reply) {
						DPrintf("[%d:%d] Get VoteReply from %d Reply Term: %d\n", rf.me, rf.currentTerm, i, reply.Term)
						if reply.Term > rf.currentTerm {
							rf.currentTerm = reply.Term
							rf.state = Follower
							return
						}
						if reply.VoteGranted {
							votes++
							DPrintf("[%d:%d] get a vote and votes is %d, need to get %d\n", rf.me, rf.currentTerm, votes, len(rf.peers)/2+1)
							if votes == len(rf.peers)/2+1 {
								rf.state = Leader
								for index := range rf.nextIndex {
									rf.nextIndex[index] = len(rf.log)
									//rf.matchIndex[index] = 1
									DPrintf("%d New nextIndex = %d", index, rf.nextIndex[index])
								}
								go rf.Heartbeat()
								return
							}
						}
					}
				}(i)
			}
			rf.electionTimer.Reset(time.Millisecond * time.Duration(rand.Int()%300+500))
		case <-rf.electionResetCh:
			DPrintf("%d reset Election Timer", rf.me)
			if !rf.electionTimer.Stop() {
				<-rf.electionTimer.C
			}
			rf.electionTimer.Reset(time.Millisecond * time.Duration(rand.Int()%300+500))
		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.applyCh = applyCh

	// Your initialization code here (2A, 2B, 2C).
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.electionResetCh = make(chan struct{}, 1)
	rf.commitCh = make(chan struct{}, 1)

	rf.log = append(rf.log, LogEntry{0, nil}) // add log 0, First log index shoubld be 1

	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
	for index := range rf.nextIndex {
		rf.nextIndex[index] = 1
		//rf.matchIndex[index] = 1
	}
	rf.lastApplied = 1

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	go rf.LeaderElection()
	go rf.CommitLogs()
	return rf
}
