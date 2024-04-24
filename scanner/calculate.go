package scanner

func CalculateReachability(urlData []URLStatus, useIPV6 bool, urlDataWithReachability []URLDataWithReachability) []URLDataWithReachability {
	children := make(map[string][]URLStatus)
	for _, data := range urlData {
		children[data.FatherURL] = append(children[data.FatherURL], data)
	}

	reachabilityIPv4 := make(map[string]URLReachabilityIPv4)
	reachabilityIPv6 := make(map[string]URLReachabilityIPv6)

	// 分别计算 IPv4 和 IPv6 的可达率
	for _, data := range urlData {
		var firstLevelReach, secondLevelTotal, secondLevelReach float64
		firstLevelChildren := children[data.URL]

		for _, child := range firstLevelChildren {
			if child.IPVersion == data.IPVersion && child.StatusCode == 200 {
				firstLevelReach++
			}

			secondLevelChildren := children[child.URL]
			for _, grandChild := range secondLevelChildren {
				if grandChild.IPVersion == data.IPVersion && grandChild.StatusCode == 200 {
					secondLevelReach++
				}
				secondLevelTotal++
			}
		}

		if len(firstLevelChildren) > 0 {
			firstLevelReach /= float64(len(firstLevelChildren))
			if useIPV6 {
				firstLevelReach *= 2
			}
		}
		if secondLevelTotal > 0 {
			secondLevelReach /= secondLevelTotal
			if useIPV6 {
				secondLevelReach *= 2
			}
		}

		if data.IPVersion == "ipv4" {
			reachabilityIPv4[data.URL] = URLReachabilityIPv4{
				FirstLevelReach:  firstLevelReach,
				SecondLevelReach: secondLevelReach,
			}
		} else if data.IPVersion == "ipv6" {
			reachabilityIPv6[data.URL] = URLReachabilityIPv6{
				FirstLevelReach:  firstLevelReach,
				SecondLevelReach: secondLevelReach,
			}
		}
	}

	for _, data := range urlData {
		newData := URLDataWithReachability{
			URLStatus:        data,
			ReachabilityIPv4: reachabilityIPv4[data.URL],
			ReachabilityIPv6: reachabilityIPv6[data.URL],
		}
		urlDataWithReachability = append(urlDataWithReachability, newData)

	}
	return urlDataWithReachability
}
