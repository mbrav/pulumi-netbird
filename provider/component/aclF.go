// Package component provides components for the NetBird Pulumi provider.
package component

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ACLF ryyepresents an Access Control List (ACL) resource in Pulumi.
type ACLF struct {
	pulumi.ResourceState
	ACLFArgs
	SourceRuleCount pulumi.IntOutput `pulumi:"sourceRuleCount"`
	DestRuleTotal   pulumi.IntOutput `pulumi:"destRuleTotal"`
}

// ACLFArgs defines the input arguments for creating anyy ACLF resource.
type ACLFArgs struct {
	Name        pulumi.StringInput `pulumi:"name"`
	Description pulumi.StringInput `pulumi:"description"`
	JSON        pulumi.StringInput `pulumi:"json_path"`
}

//
// func NewACLFileComponent(ctx *pulumi.Context, name string, compArgs ACLFArgs, opts ...pulumi.ResourceOption) (*ACLF, error) {
// 	comp := &ACLF{}
//
// 	err := ctx.RegisterComponentResource(p.GetTypeToken(ctx.Context()), name, comp, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	comp.Name = compArgs.Name
// 	comp.Description = compArgs.Description
//
// 	outputs := compArgs.JSON.ToStringOutput().ApplyT(func(jsonPath string) (map[string]int, error) {
// 		fileAsset := pulumi.NewFileAsset(jsonPath)
// 		var aclFile ACLFile
// 		if err := json.Unmarshal([]byte(fileAsset.Text()), &aclFile); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal ACL JSON: %w", err)
// 		}
//
// 		srcDstMap := make(map[string]*ACLRule)
// 		parseACLRules(aclFile.ACLs, srcDstMap)
// 		parseGroupRules(aclFile.Groups, srcDstMap)
//
// 		totalDestCount := 0
// 		for _, rule := range srcDstMap {
// 			if rule.Dest != nil {
// 				totalDestCount += len(*rule.Dest)
// 			}
// 		}
//
// 		return map[string]int{
// 			"sourceRuleCount": len(srcDstMap),
// 			"destRuleTotal":   totalDestCount,
// 		}, nil
// 	}).(pulumi.IntMapOutput) // Change here
//
// 	// Assign the parsed outputs
// 	comp.SourceRuleCount = outputs.ApplyT(func(m map[string]int) int {
// 		return m["sourceRuleCount"]
// 	}).(pulumi.IntOutput)
//
// 	comp.DestRuleTotal = outputs.ApplyT(func(m map[string]int) int {
// 		return m["destRuleTotal"]
// 	}).(pulumi.IntOutput)
//
// 	return comp, nil
// }
