#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <algorithm>
#include <set>
#include <vector>
#include <queue>

#include <fstream>
#include <string>
#include <map>
#include <sstream>

#include <Snap.h>

/* Utilities */
void WaitForInput() {
	printf("Press any key to continue ...");
	std::cin.get();
}

void printIntArray(int* arr, int n) {
	for (int i = 0; i < n; ++i) {
		printf("[%d:%d]  ", i, arr[i]);
		if (i > 0 && i % 10 == 0) {
			printf("\n");
		}
	}
	printf("\n");
}

void printDoubleMatrix(double** mat, int m, int n) {
	for (int i = 0; i < m; ++i) {
		for (int j = 0; j < n; ++j) {
			printf("%f  ", mat[i][j]);
		}
		printf("\n");
	}
}

void loadIdsFromFile(const char* fileName, std::vector<int>& idList) {
	std::ifstream inputFile(fileName);
	for (std::string line; std::getline(inputFile, line);) {
		int idint = atoi(line.c_str());
		idList.push_back(idint);
	}
	return;
}

void exportNodeList(const PNGraph& graph, std::vector<int>& idList) {
	for (TNGraph::TNodeI ni = graph->BegNI(); ni != graph->EndNI(); ni++) {
		idList.push_back(ni.GetId());
	}
}


/*  */
void getOutNeighborNodeIDs(const PNGraph& graph, int srcNodeID, std::set<int>& nodeIdSet) {
	std::queue<int> q;
	q.push(srcNodeID);
	nodeIdSet.insert(srcNodeID);
	for (int level = 0; level < 2; ++level) {
		int levelCount = q.size();
		for (int i = 0; i < levelCount; ++i) {
			int curNodeId = q.front();
			q.pop();
			// Scan neigbors;
			TNGraph::TNodeI curNode = graph->GetNI(curNodeId);
			int outDeg = curNode.GetOutDeg();
			for (int j = 0; j < outDeg; ++j) {
				int curNeighborNodeID = curNode.GetOutNId(j);
				q.push(curNeighborNodeID);
				nodeIdSet.insert(curNeighborNodeID);
			}
		}
	}
}

void getInNeighborNodeIDs(const PNGraph& graph, int destNodeID, std::set<int>& nodeIdSet) {
	std::queue<int> q;
	q.push(destNodeID);
	nodeIdSet.insert(destNodeID);
	for (int level = 0; level < 2; ++level) {
		int levelCount = q.size();
		for (int i = 0; i < levelCount; ++i) {
			int curNodeId = q.front();
			q.pop();
			// Scan neigbors;
			TNGraph::TNodeI curNode = graph->GetNI(curNodeId);
			int inDeg = curNode.GetInDeg();
			for (int j = 0; j < inDeg; ++j) {
				int curNeighborNodeID = curNode.GetInNId(j);
				q.push(curNeighborNodeID);
				nodeIdSet.insert(curNeighborNodeID);
			}
		}
	}
}

PNGraph getFourHopGraph(const PNGraph& graph, int srcId, int dstId) {
	PNGraph ret = PNGraph::New();
	int shortPath = TSnap::GetShortPath(graph, srcId, dstId, true);
	// printf("shortPath is %d\n", shortPath);
	if (shortPath > 4) {
		return ret;
	}

	std::set<int> nodeIdSet1;
	getOutNeighborNodeIDs(graph, srcId, nodeIdSet1);
	std::set<int> nodeIdSet2;
	getInNeighborNodeIDs(graph, dstId, nodeIdSet2);
	std::set<int> nodeIdSet;
	nodeIdSet.insert(nodeIdSet1.begin(), nodeIdSet1.end());
	nodeIdSet.insert(nodeIdSet2.begin(), nodeIdSet2.end());

	// Add all nodes into new graph;
	std::vector<int> nodeIdList;
	for (std::set<int>::iterator setItr = nodeIdSet.begin(); setItr != nodeIdSet.end(); setItr++) {
		ret->AddNode(*setItr);
		nodeIdList.push_back(*setItr);
	}

	// Add all edges into new graph;
	int nodeNum = nodeIdList.size();
	for (int i = 0; i < nodeNum; ++i) {
		for (int j = 0; j < nodeNum; ++j) {
			if (i == j) continue;
			if (graph->IsEdge(nodeIdList[i], nodeIdList[j], true)) {
				ret->AddEdge(nodeIdList[i], nodeIdList[j]);
			}
		}
	}
	// printf("%d nodes in sub graph\n", ret->GetNodes());
	// printf("%d edges in sub graph\n", ret->GetEdges());
	// printf("%d hops\n", TSnap::GetShortPath(ret, srcId, dstId, true));

	return ret;
}

int getNumOfIndependentPaths(const PNGraph& graph, int srcNodeID, int dstNodeID) {
	int ret = 0;
	while (true) {
		PNGraph bfsGraph = TSnap::GetBfsTree(graph, srcNodeID, true, false);
		if (!bfsGraph->IsNode(dstNodeID)) {
			return ret;
		}
		printf("%d hops\n", TSnap::GetShortPath(bfsGraph, srcNodeID, dstNodeID, true));

		// Go back from dstNode to src
		int itrNodeId = dstNodeID;
		while (itrNodeId != srcNodeID) {
			TNGraph::TNodeI curNode = bfsGraph->GetNI(itrNodeId);
			int parentNodeId = curNode.GetInNId(0);

			// Delete Edges
			// graph->DelEdge(parentNodeId, itrNodeId, true);
			// Delete Node
			if (itrNodeId != dstNodeID && itrNodeId != srcNodeID) {
				graph->DelNode(itrNodeId);
			}

			itrNodeId = parentNodeId;
		}
		++ret;
	}
}


void getDistance(const PNGraph& graph, std::vector<int> srcIds, std::vector<int> dstIds, int sampleSize, TFltPrV& ret) {
	std::random_shuffle(srcIds.begin(), srcIds.end());
	std::random_shuffle(dstIds.begin(), dstIds.end());

	int distance[20];
	for (int i = 0; i < 20; distance[i++] = 0);

	int sampleCount = 0;
	for (int i = 0; i < srcIds.size(); ++i) {
		int srcNodeId = srcIds[i];
		if (!graph->IsNode(srcNodeId)) continue;
		for (int j = 0; j < dstIds.size(); ++j) {
			int dstNodeId = dstIds[j];
			if (!graph->IsNode(dstNodeId)) continue;
			int shortDist = TSnap::GetShortPath(graph, srcNodeId, dstNodeId, true);
			distance[shortDist]++;
			sampleCount++;

			printIntArray(distance, 20);
		}
		if (sampleCount > sampleSize) break;
	}
	
	for (int i = 0; i < 20; ++i) {
		ret.Add(TFltPr(i, distance[i]));
	}
}

void getSampledDistance(const PNGraph& graph, std::vector<int> srcIds, std::vector<int> dstIds, int sampleSize, TFltPrV& ret) {
	std::random_shuffle(srcIds.begin(), srcIds.end());
	std::random_shuffle(dstIds.begin(), dstIds.end());

	int distance[20];
	for (int i = 0; i < 20; distance[i++] = 0);

	int sampleCount = 0;
	for (int i = 0; i < sampleSize; ) {
		int srcNodeId = srcIds[rand() % srcIds.size()];
		int dstNodeId = dstIds[rand() % dstIds.size()];

		if (!graph->IsNode(srcNodeId)) continue;
		if (!graph->IsNode(dstNodeId)) continue;
		int shortDist = TSnap::GetShortPath(graph, srcNodeId, dstNodeId, true);
		distance[shortDist]++;
		sampleCount++;
		printIntArray(distance, 20);
		++i;
	}

	for (int i = 0; i < 20; ++i) {
		ret.Add(TFltPr(i, distance[i]));
	}
}

void plotDegDistribution(const PNGraph& graph) {
	TFltPrV outDegDist;
	TSnap::GetOutDegCnt(graph, outDegDist);
	TGnuPlot plot1("outDegDist", "");
	plot1.AddPlot(outDegDist, gpwPoints, "");
	plot1.SetScale(gpsLog10XY);
	plot1.SavePng();

	TFltPrV inDegDist;
	TSnap::GetInDegCnt(graph, inDegDist);
	TGnuPlot plot2("inDegDist", "");
	plot2.AddPlot(inDegDist, gpwPoints, "");
	plot2.SetScale(gpsLog10XY);
	plot2.SavePng();

	TGnuPlot plot3("DegDist", "");
	plot3.AddCmd("set key right top");
	plot3.AddPlot(inDegDist, gpwPoints, "In degree");
	plot3.AddPlot(outDegDist, gpwPoints, "Out degree");
	plot3.SetScale(gpsLog10XY);
	plot3.SavePng();
}

void plotParitialDegDistribution(const PNGraph& graph, std::vector<int>& nodeList) {
	std::map<int, int> inDegDistMap;
	std::map<int, int> outDegDistMap;
	
	for (int i = 0; i < nodeList.size(); ++i) {
		int curNodeId = nodeList[i];
		if (!graph->IsNode(curNodeId)) continue;
		TNGraph::TNodeI ni = graph->GetNI(curNodeId);

		int curNodeInDeg = ni.GetInDeg();
		if (inDegDistMap.find(curNodeInDeg) == inDegDistMap.end()) {
			inDegDistMap.insert(std::pair<int, int>(curNodeInDeg, 0));
		}
		inDegDistMap[curNodeInDeg]++;

		int curNodeOutDeg = ni.GetOutDeg();
		if (outDegDistMap.find(curNodeOutDeg) == outDegDistMap.end()) {
			outDegDistMap.insert(std::pair<int, int>(curNodeOutDeg, 0));
		}
		outDegDistMap[curNodeOutDeg]++;
		
	}
	
	TFltPrV inDegDist;
	for (std::map<int, int>::iterator itr = inDegDistMap.begin(); itr != inDegDistMap.end(); itr++) {
		inDegDist.Add(TFltPr(itr->first, itr->second));
	}

	TFltPrV outDegDist;
	for (std::map<int, int>::iterator itr = outDegDistMap.begin(); itr != outDegDistMap.end(); itr++) {
		outDegDist.Add(TFltPr(itr->first, itr->second));
	}
	
	TGnuPlot plot1("inDegDistParitial", "");
	plot1.AddPlot(inDegDist, gpwPoints, "");
	plot1.SetScale(gpsLog10XY);
	plot1.SavePng();

	TGnuPlot plot2("outDegDistParitial", "");
	plot2.AddPlot(outDegDist, gpwPoints, "");
	plot2.SetScale(gpsLog10XY);
	plot2.SavePng();

	TGnuPlot plot3("DegDistParitial", "");
	plot3.AddCmd("set key right top");
	plot3.AddPlot(inDegDist, gpwPoints, "In Degree");
	plot3.AddPlot(outDegDist, gpwPoints, "Out Degree");
	plot3.SetScale(gpsLog10XY);
	plot3.SavePng();
}

void getNumOfPathsFromVect(const PNGraph& graph, std::vector<int> srcIds, int srcSampleSz, std::vector<int> dstIds, int dstSampleSz, char* fileName) {
	std::random_shuffle(srcIds.begin(), srcIds.end());
	std::random_shuffle(dstIds.begin(), dstIds.end());
	std::ofstream outputFile;
	outputFile.open(fileName);

	for (int i = 0; i < srcIds.size() && i < srcSampleSz; ++i) {
		int srcNodeId = srcIds[i];
		if (!graph->IsNode(srcNodeId)) continue;

		for (int j = 0; j < dstIds.size() && j < dstSampleSz; ++j) {
			int dstNodeId = dstIds[j];
			if (!graph->IsNode(dstNodeId)) continue;
			int shortPath = TSnap::GetShortPath(graph, srcNodeId, dstNodeId, true);
			if (shortPath > 4 || shortPath <= 2) continue;

			int numOfPaths = getNumOfIndependentPaths(graph, srcNodeId, dstNodeId);
			
			char buffer[100];
			sprintf(buffer, "%d\t%d\t%d", srcNodeId, dstNodeId, numOfPaths);
			std::cout << buffer << std::endl;
			outputFile << buffer << std::endl;
		}
	}
	outputFile.close();
}

void getPageRankFromVect(const PNGraph& graph, std::vector<int> srcIds, std::vector<int> dstIds, int sampleSz, char* fileName) {
	std::random_shuffle(srcIds.begin(), srcIds.end());
	std::random_shuffle(dstIds.begin(), dstIds.end());
	std::ofstream outputFile;
	outputFile.open(fileName);

	for (int i = 0; i < sampleSz; ) {
		int srcNodeId = srcIds[rand() % srcIds.size()];
		int dstNodeId = dstIds[rand() % dstIds.size()];
		if (!graph->IsNode(srcNodeId)) continue;
		if (!graph->IsNode(dstNodeId)) continue;
		int shortPath = TSnap::GetShortPath(graph, srcNodeId, dstNodeId, true);
		if (shortPath > 4 || shortPath <= 2) continue;

		PNGraph subgraph = getFourHopGraph(graph, srcNodeId, dstNodeId);
		TIntFltH pageRankScores;
		TSnap::GetPageRank(subgraph, pageRankScores);

		// Calculate total PR score;
		/*
		double totalPR = 0.0;
		for (TIntFltH::TIter itr = pageRankScores.BegI(); itr != pageRankScores.EndI(); itr++) {
			totalPR += itr.GetDat();
		}*/

		int numOfNodesInSubGraph = subgraph->GetNodes();
		double normalizedSrcPR = pageRankScores.GetDat(srcNodeId) * numOfNodesInSubGraph;
		double normalizedDstPR = pageRankScores.GetDat(dstNodeId) * numOfNodesInSubGraph;

		char buffer[100];
		printf("%d, %d\n", i, numOfNodesInSubGraph);
		sprintf(buffer, "%d\t%f\t%d\t%f", srcNodeId, normalizedSrcPR, dstNodeId, normalizedDstPR);
		std::cout << buffer << std::endl;
		outputFile << buffer << std::endl;
		++i;
	}
	outputFile.close();
}


void plotPR(char* fileName, TFltPrV& ret) {
	int distance[10000];
	for (int i = 0; i < 10000; distance[i++] = 0);

	std::ifstream inputFile(fileName);
	for (std::string line; std::getline(inputFile, line);) {
		std::istringstream isss(line);
		int a, c;
		double b, d;
		isss >> a >> b >> c >> d;

		int val = (int)(d * 1000);
		val -= (val % 100);
		if (val >= 10000) continue;
		//double idd = std::stold(line);
		printf("%d\n", val);
		distance[val]++;
	}

	for (int i = 0; i < 10000; ++i) {
		if (distance[i] == 0) continue;
		ret.Add(TFltPr(i, distance[i]));
	}
}

void plotpaths(char* fileName, TFltPrV& ret) {
	int distance[10000];
	for (int i = 0; i < 10000; distance[i++] = 0);

	int lineCount = 1;
	std::ifstream inputFile(fileName);
	for (std::string line; std::getline(inputFile, line);) {
		std::istringstream isss(line);
		int a, c;
		double b, d;
		isss >> a;
		ret.Add(TFltPr(lineCount++, a));
	}
}


int main(int argc, const char* argv[]) {
	// Load Twitter graph.
	PNGraph graph = TSnap::LoadEdgeList<PNGraph>("twitter_combined.txt", 0, 1);
	int numOfNodes = graph->GetNodes();
	int numOfEdges = graph->GetEdges();
	printf("Number of nodes: %d\n", numOfNodes);
	printf("Number of edges: %d\n", numOfEdges);
	
	// Load all node list.
	std::vector<int> nodeIds;
	exportNodeList(graph, nodeIds);
	
	// Load verified node list.
	std::vector<int> verifiedIds;
	loadIdsFromFile("v.txt", verifiedIds);
	printf("%d\n", verifiedIds.size());
	
	// Load spammer list.
	std::vector<int> spammerIds;
	loadIdsFromFile("labeled_spammers.txt", spammerIds);
	printf("%d\n", spammerIds.size());

	// Load suspect list.
	std::vector<int> suspectIds;
	loadIdsFromFile("suspect.txt", suspectIds);
	std::random_shuffle(suspectIds.begin(), suspectIds.end());
	printf("%d\n", suspectIds.size());


	getPageRankFromVect(graph, verifiedIds, spammerIds, 1000, "pr_VS2.txt");
	getPageRankFromVect(graph, verifiedIds, verifiedIds, 1000, "pr_VV2.txt");

	/*
	TFltPrV vvPR;
	plotPR("pr_VV2.txt", vvPR);
	TFltPrV vsPR;
	plotPR("pr_VS2.txt", vsPR);
	TGnuPlot plot3("prcom2", "");
	plot3.AddCmd("set key right top");
	plot3.SetXRange(0, 2500);
	plot3.AddPlot(vvPR, gpwLinesPoints, "Legitimate users to legitimate users");
	plot3.AddPlot(vsPR, gpwLinesPoints, "Legitimate users to Spammers");
	plot3.SetXLabel("Page rank score x 1000");
	plot3.SetYLabel("Count");
	// plot3.SetScale(gpsLog10XY);
	plot3.SavePng();
	*/
	/*
	TFltPrV vvPR;
	plotpaths("numOfPaths_VV.txt", vvPR);
	TFltPrV vsPR;
	plotpaths("numOfPaths_VS.txt", vsPR);
	TGnuPlot plot3("pathcom", "");
	plot3.AddCmd("set key right top");
	plot3.SetXRange(0, 100);
	plot3.AddPlot(vvPR, gpwLinesPoints, "Legitimate users to legitimate users");
	plot3.AddPlot(vsPR, gpwLinesPoints, "Legitimate users to Spammers");
	plot3.SetXLabel("Number of indepentdent paths");
	plot3.SetYLabel("Count");
	// plot3.SetScale(gpsLog10XY);
	plot3.SavePng();
	*/
	/*
	TFltPrV regularDistance;
	getSampledDistance(graph, verifiedIds, verifiedIds, 5000, regularDistance);
	TFltPrV spamDistance;
	getSampledDistance(graph, verifiedIds, spammerIds, 5000, spamDistance);

	TGnuPlot plot3("Distance", "");
	plot3.AddCmd("set key right top");
	plot3.SetXRange(1, 9);
	plot3.AddPlot(regularDistance, gpwLinesPoints, "Legitimate users to legitimate users");
	plot3.AddPlot(spamDistance, gpwLinesPoints, "Legitimate users to Spammers");
	plot3.SetXLabel("Distance");
	plot3.SetYLabel("Count");
	// plot3.SetScale(gpsLog10XY);
	plot3.SavePng();
	*/



	// Distance Analysis
	// plotDegDistribution(graph);
	// plotParitialDegDistribution(graph, suspectIds);
	// getDistance(graph, verifiedIds, spammerIds);

	//getPageRankFromVect(graph, verifiedIds, 30, verifiedIds, 30, "pageRankResult_VV.txt");
	//getPageRankFromVect(graph, verifiedIds, 30, spammerIds, 30, "pageRankResult_VS.txt");

	// getNumOfPathsFromVect(graph, verifiedIds, 30, verifiedIds, 30, "numOfPaths_VV.txt");
	// getNumOfPathsFromVect(graph, verifiedIds, 30, spammerIds, 30, "numOfPaths_VS.txt");
	/*

	// double aveClusterCoeff = TSnap::GetClustCf(graph);
	// printf("%f\n", aveClusterCoeff);

	int srcNodeID = 783214; // graph->GetRndNId();
	int dstNodeID = 657863; // graph->GetRndNId();
	printf("src: %d, dst: %d\n", srcNodeID, dstNodeID);

	PNGraph subgraph = getFourHopGraph(graph, srcNodeID, dstNodeID);

	TIntFltH pageRanks;
	TSnap::GetPageRank(subgraph, pageRanks);

	pageRanks.Sort(false, false);
	TVec<TInt, int> keyV;
	pageRanks.GetKeyV(keyV);

	int ccc = 0;
	for (TVec<TInt, int>::TIter itr = keyV.BegI(); itr != keyV.EndI(); itr++) {
		printf("%d, %f\n", itr->Val, pageRanks.GetDat(itr->Val).Val);
		if (++ccc % 50 == 0) {
			WaitForInput();
		}
	}
	*/

	/*
	int numPaths = getNumOfIndependentPaths(subgraph, srcNodeID, dstNodeID);
	printf("Num Of paths: %d", numPaths);
	*/
	WaitForInput();
	return 0;
}
