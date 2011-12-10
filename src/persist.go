package movo

import (
    "appengine"
    "appengine/datastore"
    "appengine/memcache"
    "bytes"
    "gob"
    "os"
    "strconv"
    )

type DSActDetail struct {
  ActId int
  Type string
  Name string
  Desc string `datastore:",noindex"`
  Location string
  Start datastore.Time
  End datastore.Time
  NCadets int `datastore:",noindex"`
  NStaff int `datastore:",noindex"`
}

func (a *ActDetail) toDS() (res *DSActDetail) {
  res = &DSActDetail{
          ActId: a.ActId,
          Type: a.Type,
          Name: a.Name,
          Desc: a.Desc,
          Location: a.Location,
          NCadets: a.NCadets,
          NStaff: a.NStaff,
        }
  if a.Start != nil {
    res.Start = datastore.SecondsToTime(a.Start.Seconds())
  }
  if a.End != nil {
    res.End = datastore.SecondsToTime(a.End.Seconds())
  }

  return
}

func (d *DSActDetail) fromDS() (res *ActDetail) {
  res = &ActDetail{
          ActId: d.ActId,
          Type: d.Type,
          Name: d.Name,
          Desc: d.Desc,
          Location: d.Location,
          Start: d.Start.Time(),
          End: d.End.Time(),
          NCadets: d.NCadets,
          NStaff: d.NStaff,
        }

  return
}

func ActExists(c appengine.Context, actId int) (bool, os.Error) {
  return false, nil
}

func RetrieveActDetails(c appengine.Context, actId int) (res *ActDetail, err os.Error) {
  var d DSActDetail
  if itm, err := memcache.Get(c, "actId__" + strconv.Itoa(actId));
      err != nil && err != memcache.ErrCacheMiss {
    return nil, err
  } else if err == nil {
    // Cache hit
    buf := bytes.NewBuffer(itm.Value)
    dec := gob.NewDecoder(buf)
    dec.Decode(&d)
  } else {
    // Cache miss
    key := datastore.NewKey(c, "DSActDetail", "", int64(actId), nil)
    if err := datastore.Get(c, key, &d); err != nil {
      return nil, err
    }
  }
  return d.fromDS(), nil
}

func (a *ActDetail) Persist(c appengine.Context) os.Error {
  d := a.toDS()
c.Debugf("Done converting to DSActDetail")
  key := datastore.NewKey(c, "DSActDetail", "", int64(d.ActId), nil)
c.Debugf("Done Creating Key")
  _,err := datastore.Put(c, key, d)
c.Debugf("Done persisting to datastore")
  if err != nil {
    return err
  }

  buf := bytes.NewBufferString("")
  enc := gob.NewEncoder(buf)
  enc.Encode(d)

  itm := &memcache.Item{
          Key:    "actId__" + strconv.Itoa(a.ActId),
          Value:  buf.Bytes(),
       }

  err = memcache.Set(c, itm)
c.Debugf("Request cache to memcache")

  return err
}

func init() {
  gob.Register(ActDetail{})
}
