import pathlib,json,hashlib,types,io,subprocess,re
from fractions import Fraction
r=pathlib.Path('/private/tmp/data-fabric-w1-p1-wire-semantic-clarification'); o=pathlib.Path(__file__).parent
m=types.ModuleType('review_importer'); m.__file__=str(r/'sources/import_ror.py'); exec(compile((r/'sources/import_ror.py').read_bytes(),m.__file__,'exec'),m.__dict__)
v=json.loads((r/'test-vectors.json').read_text())['vectors']
extra=['1E+000','1000.000','0.1e1','9007199254740991.0','9007199254740991.1','9007199254740991.0000000000000001','-9007199254740992','-1e400','1e400','0e400','false','[1,{"x":9007199254740993}]']
results=[]
def classify(t):
 if t=='null':return 'null-reference'
 if t=='""':return 'empty-reference'
 if not re.fullmatch(r'-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?',t):return 'invalid-reference'
 n=Fraction(t)
 if n<=0 or n.denominator!=1:return 'invalid-reference'
 if n>9007199254740991:return 'unsupported-precision'
 if not re.fullmatch('[1-9][0-9]*',t):return 'invalid-reference'
 return 'eligible'
for t in [x['token'] for x in v]+extra:
 doc='[{"id":"https://ror.org/000025p04","status":"active","locations":[{"geonames_id":'+t+',"geonames_details":{}}],"relationships":[]}]'
 try:
  rec=list(m.records(io.StringIO(doc)))[0]; m.validate_record(rec); native=True; err=None
 except (ValueError,OverflowError) as e:native=False;err=str(e)
 result={'token':t,'exact_classification':classify(t),'publisher_accepts':native,'publisher_error':err}
 expected=next((x for x in v if x['token']==t),None)
 if expected:
  assert result['exact_classification']==expected['expected'],t
  assert native==expected['publisher_accepts'],t
 results.append(result)
js='const ts='+json.dumps([x['token'] for x in results])+'; console.log(JSON.stringify(ts.map(token=>{let v=JSON.parse(token);return {token,number:String(v),safe:Number.isSafeInteger(v),finite:typeof v=== "number"?Number.isFinite(v):null,negativeZero:Object.is(v,-0)}})))'
js_results=json.loads(subprocess.check_output(['node','-e',js]))
malformed=['01','+1','NaN','Infinity','{"geonames_id":1,"geonames_id":2}','{"geonames_id":1,"geonames_\\u0069d":2}','{"geonames_id":']
# JSON duplicate names are an integrity policy, not a JSON grammar violation.
for t in malformed:
 try:json.loads(t,object_pairs_hook=m.unique_object,parse_constant=m.reject_constant); raise AssertionError(t)
 except (ValueError,json.JSONDecodeError):pass
report={'scope':'independent exact Fraction oracle and pinned importer controls; no product acceptance','original_vectors':len(v),'extra_controls':len(extra),'numeric_results':results,'js_results':js_results,'malformed_or_ambiguous':malformed,'duplicate_escaped_name_rejected':True,'script_sha256':hashlib.sha256(pathlib.Path(__file__).read_bytes()).hexdigest()}
(o/'probe-results.json').write_text(json.dumps(report,indent=2)+'\n'); print(json.dumps({'original_vectors':len(v),'extra_controls':len(extra),'all_matched':True,'malformed_or_ambiguous_rejected':len(malformed)}))
